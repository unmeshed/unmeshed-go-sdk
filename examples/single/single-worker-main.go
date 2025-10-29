package main

import (
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	apis "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/main"
	apis2 "github.com/unmeshed/unmeshed-go-sdk/sdk/apis/workers"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/common"
	"github.com/unmeshed/unmeshed-go-sdk/sdk/configs"
)

var unmeshedClient *apis.UnmeshedClient

func StringPtr(s string) *string                              { return &s }
func IntPtr(i int) *int                                       { return &i }
func StepTypePtr(t common.StepType) *common.StepType          { return &t }
func ProcessTypePtr(t common.ProcessType) *common.ProcessType { return &t }

func TestWorker(data map[string]interface{}) string {
	fmt.Println("TestWorker running with data:", data)
	time.Sleep(1 * time.Second)
	return "Hello from test-worker!"
}

func createProcessDefinitionExample(client *apis.UnmeshedClient) {
	log.Println("========== Creating Process Definition ==========")

	noop1Step := &common.StepDefinition{
		Name:        StringPtr("noop1"),
		Ref:         StringPtr("noop1"),
		Description: StringPtr("Test noop 1"),
		Type:        StepTypePtr(common.StepTypeNoop),
		Input:       map[string]interface{}{"key1": "val1"},
	}

	noop2Step := &common.StepDefinition{
		Name:        StringPtr("noop2"),
		Ref:         StringPtr("noop2"),
		Description: StringPtr("Test noop 2"),
		Type:        StepTypePtr(common.StepTypeNoop),
		Input:       map[string]interface{}{"key1": "val1"},
	}

	processDefinition := &common.ProcessDefinition{
		Name:        StringPtr("test-process"),
		Version:     IntPtr(1),
		Namespace:   StringPtr("default"),
		Description: StringPtr("Testing Process"),
		Type:        ProcessTypePtr(common.ProcessTypeAPIOrchestration),
		Steps:       []*common.StepDefinition{noop1Step, noop2Step},
	}

	createdPD, err := client.CreateNewProcessDefinition(processDefinition)
	if err != nil {
		log.Printf("Error creating process definition: %v\n", err)
		return
	}

	log.Printf("Created process definition: %+v\n", createdPD)
	log.Printf("Created process definition name: %s, version: %d\n",
		*createdPD.Name, *createdPD.Version)
}

func updateProcessDefinitionExample(client *apis.UnmeshedClient) {
	log.Println("\n========== Updating Process Definition ==========")

	noop1Step := &common.StepDefinition{
		Name:        StringPtr("noop1"),
		Ref:         StringPtr("noop1"),
		Description: StringPtr("Test noop 1"),
		Type:        StepTypePtr(common.StepTypeNoop),
		Input:       map[string]interface{}{"key1": "val1"},
	}

	noop2Step := &common.StepDefinition{
		Name:        StringPtr("noop2"),
		Ref:         StringPtr("noop2"),
		Description: StringPtr("Test noop 2"),
		Type:        StepTypePtr(common.StepTypeNoop),
		Input:       map[string]interface{}{"key1": "val1"},
	}

	noop3Step := &common.StepDefinition{
		Name:        StringPtr("noop3"),
		Ref:         StringPtr("noop3"),
		Description: StringPtr("Test noop 3"),
		Type:        StepTypePtr(common.StepTypeNoop),
		Input:       map[string]interface{}{"key1": "val1"},
	}

	processDefinition := &common.ProcessDefinition{
		Name:        StringPtr("test-process"),
		Version:     IntPtr(2),
		Namespace:   StringPtr("default"),
		Description: StringPtr("Testing Process Updated"),
		Type:        ProcessTypePtr(common.ProcessTypeAPIOrchestration),
		Steps:       []*common.StepDefinition{noop1Step, noop2Step, noop3Step},
	}

	updatedPD, err := client.UpdateProcessDefinition(processDefinition)
	if err != nil {
		log.Printf("Error updating process definition: %v\n", err)
		return
	}

	log.Printf("Updated process definition: %+v\n", updatedPD)
	log.Printf("Updated process definition name: %s, version: %d, description: %s\n",
		*updatedPD.Name, *updatedPD.Version, *updatedPD.Description)
}

func getProcessDefinitionExample(client *apis.UnmeshedClient) {
	log.Println("\n========== Getting Process Definition (Latest) ==========")

	latestPD, err := client.GetProcessDefinitionLatestOrVersion("default", "test-process", nil)
	if err != nil {
		log.Printf("Error getting latest process definition: %v\n", err)
		return
	}

	log.Printf("Latest process definition: name=%s, version=%d, description=%s\n",
		*latestPD.Name, *latestPD.Version, *latestPD.Description)
	log.Printf("Steps count: %d\n", len(latestPD.Steps))

	version := 1
	log.Println("\n========== Getting Process Definition (Specific Version) ==========")
	specificPD, err := client.GetProcessDefinitionLatestOrVersion("default", "test-process", &version)
	if err != nil {
		log.Printf("Error getting specific process definition version: %v\n", err)
		return
	}

	log.Printf("Process definition v%d: name=%s, description=%s\n",
		*specificPD.Version, *specificPD.Name, *specificPD.Description)
	log.Printf("Steps count: %d\n", len(specificPD.Steps))
}

func getProcessDefinitionVersionsExample(client *apis.UnmeshedClient) {
	log.Println("\n========== Getting Process Definition Versions ==========")

	versions, err := client.GetProcessDefinitionVersions("default", "test-process")
	if err != nil {
		log.Printf("Error getting process definition versions: %v\n", err)
		return
	}

	log.Printf("Available versions for test-process: %v\n", versions)
	log.Printf("Total versions available: %d\n", len(versions))
}

func getAllProcessDefinitionsExample(client *apis.UnmeshedClient) {
	log.Println("\n========== Getting All Process Definitions ==========")

	allPDs, err := client.GetAllProcessDefinitions()
	if err != nil {
		log.Printf("Error getting all process definitions: %v\n", err)
		return
	}

	log.Printf("Total process definitions: %d\n", len(allPDs))

	var testProcessDefs []*common.ProcessDefinition
	for _, pd := range allPDs {
		if *pd.Name == "test-process" && *pd.Namespace == "default" {
			testProcessDefs = append(testProcessDefs, pd)
		}
	}

	log.Printf("Filtered test-process definitions: %d\n", len(testProcessDefs))
	for _, pd := range testProcessDefs {
		log.Printf("  - %s (v%d): %s\n", *pd.Name, *pd.Version, *pd.Description)
	}
}

func deleteProcessDefinitionsExample(client *apis.UnmeshedClient) {
	log.Println("\n========== Deleting Process Definitions ==========")

	allPDs, err := client.GetAllProcessDefinitions()
	if err != nil {
		log.Printf("Error getting all process definitions: %v\n", err)
		return
	}

	var testProcessDefs []*common.ProcessDefinition
	for _, pd := range allPDs {
		if *pd.Name == "test-process" && *pd.Namespace == "default" {
			testProcessDefs = append(testProcessDefs, pd)
		}
	}

	if len(testProcessDefs) > 0 {
		log.Printf("Deleting %d test-process definitions...\n", len(testProcessDefs))

		deleteResponse, err := client.DeleteProcessDefinitions(testProcessDefs, false)
		if err != nil {
			log.Printf("Error deleting process definitions: %v\n", err)
			return
		}

		log.Printf("Delete response: %+v\n", deleteResponse)
		log.Println("Process definitions deleted successfully")
	} else {
		log.Println("No test-process definitions found to delete")
	}
}

func deleteSpecificVersionExample(client *apis.UnmeshedClient) {
	log.Println("\n========== Deleting Specific Version ==========")

	allPDs, err := client.GetAllProcessDefinitions()
	if err != nil {
		log.Printf("Error getting all process definitions: %v\n", err)
		return
	}

	var version1Defs []*common.ProcessDefinition
	for _, pd := range allPDs {
		if *pd.Name == "test-process" && *pd.Namespace == "default" && *pd.Version == 1 {
			version1Defs = append(version1Defs, pd)
		}
	}

	if len(version1Defs) > 0 {
		log.Printf("Deleting version 1 of test-process...\n")

		deleteResponse, err := client.DeleteProcessDefinitions(version1Defs, true)
		if err != nil {
			log.Printf("Error deleting specific version: %v\n", err)
			return
		}

		log.Printf("Delete response: %+v\n", deleteResponse)
		log.Println("Version 1 deleted successfully")
	}
}

func main() {
	worker := apis2.NewWorker(TestWorker, "test-worker")

	clientConfig := configs.NewClientConfig()
	clientConfig.SetClientID("<< Client Id >>")
	clientConfig.SetAuthToken("<< Auth Token>>")
	clientConfig.SetDelayMillis(40)
	clientConfig.SetPort(8080)
	clientConfig.SetWorkRequestBatchSize(500)
	clientConfig.SetBaseURL("http://localhost")
	clientConfig.SetStepTimeoutMillis(36000000)
	clientConfig.SetMaxWorkers(100)

	var err error
	unmeshedClient, err = apis.NewUnmeshedClient(clientConfig)
	if err != nil {
		fmt.Printf("Error creating client: %v\n", err)
		return
	}

	createProcessDefinitionExample(unmeshedClient)
	updateProcessDefinitionExample(unmeshedClient)
	getProcessDefinitionExample(unmeshedClient)
	getProcessDefinitionVersionsExample(unmeshedClient)
	getAllProcessDefinitionsExample(unmeshedClient)
	deleteProcessDefinitionsExample(unmeshedClient)

	unmeshedClient.RegisterWorker(worker)

	done := make(chan struct{})
	go func() {
		unmeshedClient.Start()
		close(done)
	}()

	fmt.Println("Single worker client started. Press Ctrl+C to stop...")

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, os.Interrupt, syscall.SIGTERM)

	select {
	case <-sigChan:
		fmt.Println("\nReceived shutdown signal. Stopping client...")
	case <-done:
		fmt.Println("Client finished execution")
	}
}
