package main

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/ahmetb/go-linq"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/lambda"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/montanaflynn/stats"
)

func upload() {
	// generate 16 objects with size range from 1KB to 32 MB, increase by a factor of 2
	initSizeInBytes := 1024
	var i uint8
	for i = 1; i <= 16; i++ {
		subKey := getObjectName(i)
		_, err := g_s3_service.HeadObject(&s3.HeadObjectInput{
			Bucket: aws.String(g_bucket_name),
			Key:    aws.String(subKey),
		})

		if err != nil {
			//could not find object
			input := &s3.PutObjectInput{
				Body:   bytes.NewReader(make([]byte, initSizeInBytes)),
				Bucket: aws.String(g_bucket_name),
				Key:    aws.String(subKey),
			}
			_, err := g_s3_service.PutObject(input)

			if err != nil {
				recordError(err)
			}
		}
		initSizeInBytes *= 2
	}
}

func getFunctionConfigByName(functionName string) *lambda.FunctionConfiguration {
	input := &lambda.GetFunctionConfigurationInput{
		FunctionName: aws.String(functionName),
	}

	result, err := g_lambda_service.GetFunctionConfiguration(input)
	if err != nil {
		recordError(err)
	}
	return result
}

type ReportInfo struct {
	ProfileName      string
	MemorySizeInMB   int64
	ConcurrentNumber int
	RawJson          string
}

type ReportFiles struct {
	RawReport   string
	StatsReport string
	ProfileName string
}

func generate_report(prefix []byte, info ReportInfo) ReportFiles {
	finalReportFiles := ReportFiles{}
	// get report units from S3
	prefixInStr := strings.Trim(string(prefix[:]), "\"")
	fmt.Printf("get report units from S3, key is %s ...\n", prefixInStr)
	report_units := downloadByPrefix(g_bucket_name, prefixInStr)
	if len(report_units) == 0 {
		return finalReportFiles
	}

	// parse headers
	fmt.Println("parse headers ...")
	headersToMap := map[string][]float64{}
	headers := []string{}
	m := map[string][]float64{}
	err := json.Unmarshal(report_units[0], &m)
	if err != nil {
		recordError(err)
		return finalReportFiles
	}
	for key, _ := range m {
		headersToMap[key] = []float64{}
		headers = append(headers, key)
	}

	// aggregate all report units
	fmt.Println("aggregate all report units ...")
	for _, item := range report_units {
		m = map[string][]float64{}
		err := json.Unmarshal(item, &m)
		if err != nil {
			recordError(err)
			continue
		}

		for key, objs := range m {
			headersToMap[key] = append(headersToMap[key], objs...)
		}
	}

	// sort each metrics in descending order
	fmt.Println("sort each metrics in descending order ...")
	for _, key := range headers {
		sort.Sort(sort.Reverse(sort.Float64Slice(headersToMap[key])))
	}

	// do stats such as mean, p99, min etc.
	fmt.Println("do stats such as mean, p99, min etc. ...")
	flat_data := []float64{}
	record_number := len(headersToMap[headers[0]])
	headers_number := len(headers)
	var statBuffer strings.Builder
	testConditions := fmt.Sprintf("Memory size: %d, Concurrent number: %d", info.MemorySizeInMB, info.ConcurrentNumber)
	rawjsonObj := map[string]interface{}{}
	err = json.Unmarshal([]byte(info.RawJson), &rawjsonObj)
	if err == nil {
		for key, val := range rawjsonObj {
			s, ok := val.(string)
			if ok {
				testConditions = testConditions + fmt.Sprintf(", %s: %s", key, s)
				continue
			}
			d, ok := val.(float64)
			if ok {
				testConditions = testConditions + fmt.Sprintf(", %s: %f", key, d)
				continue
			}
		}
	}
	statBuffer.WriteString(testConditions)
	statBuffer.WriteString("\n")
	headerCols := strings.Join(headers[:], ",")
	statBuffer.WriteString(fmt.Sprintf("Stats Metrics,%s\n", headerCols))

	metrics := [][]string{[]string{"avg"}, []string{"min"}, []string{"p25"}, []string{"p50"}, []string{"p75"}, []string{"p90"}, []string{"p99"}, []string{"max"}}
	for _, key := range headers {
		avg, _ := stats.Mean(headersToMap[key])
		metrics[0] = append(metrics[0], fmt.Sprintf("%f", avg))
		min := headersToMap[key][len(headersToMap[key])-1]
		metrics[1] = append(metrics[1], fmt.Sprintf("%f", min))
		p25, _ := stats.Percentile(headersToMap[key], 25)
		metrics[2] = append(metrics[2], fmt.Sprintf("%f", p25))
		p50, _ := stats.Percentile(headersToMap[key], 50)
		metrics[3] = append(metrics[3], fmt.Sprintf("%f", p50))
		p75, _ := stats.Percentile(headersToMap[key], 75)
		metrics[4] = append(metrics[4], fmt.Sprintf("%f", p75))
		p90, _ := stats.Percentile(headersToMap[key], 90)
		metrics[5] = append(metrics[5], fmt.Sprintf("%f", p90))
		p99, _ := stats.Percentile(headersToMap[key], 99)
		metrics[6] = append(metrics[6], fmt.Sprintf("%f", p99))
		max := headersToMap[key][0]
		metrics[7] = append(metrics[7], fmt.Sprintf("%f", max))

		flat_data = append(flat_data, headersToMap[key]...)
	}

	for _, metric := range metrics {
		statBuffer.WriteString(strings.Join(metric[:], ","))
		statBuffer.WriteString("\n")
	}

	// generate report
	fmt.Println("generate report ...")
	var buffer strings.Builder
	buffer.WriteString(testConditions)
	buffer.WriteString("\n")
	buffer.WriteString(headerCols)
	buffer.WriteString("\n")
	for i := 0; i < record_number; i++ {
		one_row := []float64{}
		for j := 0; j < headers_number; j++ {
			one_row = append(one_row, flat_data[i+j*record_number])
		}
		buffer.WriteString(strings.ReplaceAll(strings.Trim(fmt.Sprint(one_row), "[]"), " ", ","))
		buffer.WriteString("\n")
	}

	t := time.Now()
	dt := fmt.Sprintf("%d-%02d-%02dT%02d_%02d_%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	d1 := []byte(strings.Trim(buffer.String(), "\n"))
	finalReportFiles.RawReport = getReportPath(strings.TrimSpace(fmt.Sprintf("raw-data-%s-%s-%s.csv", info.ProfileName, dt, prefixInStr)))
	err = ioutil.WriteFile(finalReportFiles.RawReport, d1, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
	d2 := []byte(strings.Trim(statBuffer.String(), "\n"))
	finalReportFiles.StatsReport = getReportPath(strings.TrimSpace(fmt.Sprintf("report-%s-%s-%s.csv", info.ProfileName, dt, prefixInStr)))
	err = ioutil.WriteFile(finalReportFiles.StatsReport, d2, 0644)
	if err != nil {
		fmt.Println(err.Error())
	}
	finalReportFiles.ProfileName = info.ProfileName
	return finalReportFiles
}

func getReportPath(fileName string) string {
	return fmt.Sprintf("reports/%s", fileName)
}

func mergeReports(reports []interface{}) {
	t := time.Now()
	d := fmt.Sprintf("%d-%02d-%02dT%02d_%02d_%02d",
		t.Year(), t.Month(), t.Day(),
		t.Hour(), t.Minute(), t.Second())
	{
		var buffer strings.Builder
		for _, r := range reports {
			content, _ := ioutil.ReadFile(r.(ReportFiles).RawReport)
			// Convert []byte to string and write to file
			buffer.WriteString(string(content))
			buffer.WriteString("\n")
		}
		ioutil.WriteFile(getReportPath(fmt.Sprintf("raw-data-%s-%s.csv", reports[0].(ReportFiles).ProfileName, d)), []byte(strings.Trim(buffer.String(), "\n")), 0644)
	}
	{
		var buffer strings.Builder
		for _, r := range reports {
			content, _ := ioutil.ReadFile(r.(ReportFiles).StatsReport)
			// Convert []byte to string and write to file
			buffer.WriteString(string(content))
			buffer.WriteString("\n")
		}
		ioutil.WriteFile(getReportPath(fmt.Sprintf("report-%s-%s.csv", reports[0].(ReportFiles).ProfileName, d)), []byte(strings.Trim(buffer.String(), "\n")), 0644)
	}
}

var g_bucket_name string

func main() {
	timeToWaitArg := flag.Int("time-to-wait", 1, "Time to wait when begins to get reports in S3, unit by Minute.")
	bucketNameArg := flag.String("bucket-name", "", "Bucket name to store generated reports.")
	testDeploymentArg := flag.Bool("test-deployment", false, "User it to test whether the infrastructures are deployed properly, set true to validate.")
	flag.Parse()
	if len(*bucketNameArg) == 0 {
		fmt.Println("Please provide bucket name, for example, enter the following command:")
		fmt.Println("./auto-run <your-bucket-name>")
		return
	}

	g_bucket_name = *bucketNameArg
	init_shared_resource()
	// launch Lambda Function
	var params []EventParams
	if *testDeploymentArg {
		params = []EventParams{
			EventParams{NumberOfTasks: 6, LambdaFunctionName: "worker-handler", TaskName: "DefaultPerformancer"},
		}
	} else {
		// Open our jsonFile
		jsonFile, err := os.Open("config.json")
		// if we os.Open returns an error then handle it
		if err != nil {
			fmt.Println(err)
			return
		}
		// read our opened jsonFile as a byte array.
		byteValue, err := ioutil.ReadAll(jsonFile)
		if err != nil {
			fmt.Println(err)
			return
		}

		// we unmarshal our byteArray which contains our
		// jsonFile's content into 'users' which we defined above
		json.Unmarshal(byteValue, &params)
		for i := 0; i < len(params); i++ {
			params[i].LambdaFunctionName = "worker-handler"
		}
	}

	fmt.Println("Start ...")
	results := [][]byte{}
	for _, p := range params {
		payLoadInJson, _ := json.Marshal(p)
		input := &lambda.InvokeInput{
			FunctionName: aws.String("test-harness-framework"),
			Payload:      payLoadInJson,
		}
		result, err := g_lambda_service.Invoke(input)
		if err != nil {
			recordError(err)
			continue
		}
		fmt.Println(fmt.Sprintf("Task %s is launched", string(result.Payload[:])))
		results = append(results, result.Payload)
	}

	// wait after timeToWaitArg minutes to begin collect reports
	fmt.Println(fmt.Sprintf("Start to wait about %d minutes ...", *timeToWaitArg))
	time.Sleep(time.Duration(*timeToWaitArg) * time.Minute)

	// generate reports
	reports := []ReportFiles{}
	for idx, item := range results {
		fc := getFunctionConfigByName(params[idx].LambdaFunctionName)
		info := ReportInfo{ProfileName: params[idx].TaskName, MemorySizeInMB: *fc.MemorySize, ConcurrentNumber: params[idx].ConcurrencyForEachTask, RawJson: params[idx].RawJson}
		reports = append(reports, generate_report(item, info))
	}

	// merge reports
	fmt.Println("Merge reports ...")
	q := linq.From(reports).GroupBy(
		func(i interface{}) interface{} { return i.(ReportFiles).ProfileName },
		func(i interface{}) interface{} { return i.(ReportFiles) })

	groupedResults := q.OrderBy(func(i interface{}) interface{} {
		return i.(linq.Group).Key
	}).Results()
	for _, r := range groupedResults {
		mergeReports(r.(linq.Group).Group)
	}
	fmt.Println("End!")
}
