package main

import "os"
import "rialto"
import "fmt"
//import "encoding/json"
import "io/ioutil"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func checkOk(ok bool, message string) {
	if !ok {
		errStr:=fmt.Sprintf(message)
		//fmt.Println(errStr)
		panic(errStr)
	}
}

func main() {
    /*
    var wellKnownTypesFile string = "../test/well-known-types.json"
    rialto.ReadWellKnownServiceConfig( wellKnownTypesFile )
    */


	//Initialize, load available services into memory
	//TODO:err
	context:=rialto.InitMe()

	//TODO: Path should be a parameter
    	var installConfigFile string = "../test/install-env.json"

	//parse config
    	installConfig:= rialto.ReadInstallConfig( installConfigFile )

	//setup a list of running service instances
    	serviceInstances:=map[string]rialto.ServiceInstance{}

	//setup deployment environement
	var env rialto.EnvironmentInstance = rialto.EnvironmentInstance{}
	env.ExternalIPs=[]string{"asdas", "asdsadasdasd"}

	f, err := os.Create("./deploy.sh")
	check(err)
	defer f.Close()
	f.WriteString("#!/bin/bash\n\n")
	for _,plannedServerInstances := range installConfig.Instances {
		foundServiceConfig, ok:= context.Services[plannedServerInstances.ServiceName]
		checkOk(ok, fmt.Sprintf("Service not found: %v", plannedServerInstances.Name))

		//fmt.Println(foundServiceConfig.Name)

        	//create service instance from config
        	serviceInstances[plannedServerInstances.Name]=rialto.GetServiceInstance(foundServiceConfig, plannedServerInstances.Name )

		fmt.Println(serviceInstances[plannedServerInstances.Name])

	    	// load values file
		var deploymentEnvironment rialto.Deployment= rialto.Deployment{}
		deploymentEnvironment.ServiceInstance=serviceInstances[plannedServerInstances.Name]
		deploymentEnvironment.Environment=env

		//write our the values file
		var templateStr string = rialto.LoadValueTemplate(plannedServerInstances)
		//fmt.Println(templateStr)
		valueBytes:=rialto.ApplyTemplate(templateStr, deploymentEnvironment)
		var valuesFileName string =  fmt.Sprintf("./values.%v.yaml", plannedServerInstances.Name)
		err:= ioutil.WriteFile(valuesFileName, valueBytes, 0644)
	    	check(err)//fmt.Println(valueBuf.String())

		f.WriteString(fmt.Sprintf("helm install %v\n", plannedServerInstances.ChartSource.HelmRepo))





	/*
        configFile := map[string]string{}
        for _,interfaceVal := range foundServiceConfig.Exposes {
            for _,propVal := range interfaceVal.Properties {
                configFile[propVal.Name]="???"
            }
        }
	*/
    }

//    fmt.Println(configFiles)
}
