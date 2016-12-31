package main

import "os"
import "rialto"
import "fmt"
//import "encoding/json"
import "io/ioutil"
import "flag"

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
	//define a flag - use: rialto -install-file=../test/install.json
	//defaults to ./install.json
	installPtr := flag.String("install-file", "./install.json", "Installation specification file(see ???)")
	debugPrt := flag.Bool("debug", false, "Verbose if set")
	//parse flag and set variables
	flag.Parse()
	installConfigFile:= *installPtr
	debug:=*debugPrt

	if(debug) { fmt.Println("Using install file:", installConfigFile) }

	//Initialize, load available services into memory
	//TODO:err
	//context:=rialto.Init()

	//TODO: Path should be a parameter
    	//var installConfigFile string = "../test/install-env.json"

	//parse config
    	installConfig:= rialto.ReadInstallConfig( installConfigFile )

	//setup a list of running service instances
    	serviceInstances:=map[string]rialto.ServiceInstance{}

	//setup deployment environement
	var env rialto.EnvironmentInstance = rialto.EnvironmentInstance{}
	env.ExternalIPs=rialto.GetExternalIPs();

	//env.ExternalIPs=[]string{"asdas", "asdsadasdasd"}

	f, err := os.Create("./deploy.sh")
	check(err)
	defer f.Close()
	f.WriteString("#!/bin/bash\n\n")
	for _,plannedServiceInstance := range installConfig.Instances {
		foundServiceConfig/*, ok*/:= rialto.LoadServiceConfig(plannedServiceInstance)// context.Services[plannedServiceInstance.ServiceName]
		/*checkOk(ok, fmt.Sprintf("Service not found: %v", plannedServiceInstance.Name))
*/
		//fmt.Println(foundServiceConfig.Name)

        	//create service instance from config
        	serviceInstances[plannedServiceInstance.Name]=rialto.GetServiceInstance(foundServiceConfig, plannedServiceInstance.Name )

		fmt.Println(serviceInstances[plannedServiceInstance.Name])

	    	// load values file
		var deploymentEnvironment rialto.Deployment= rialto.Deployment{}
		deploymentEnvironment.ServiceInstance=serviceInstances[plannedServiceInstance.Name]
		deploymentEnvironment.Environment=env

		//write our the values file
		var templateStr string = rialto.LoadValueTemplate(plannedServiceInstance)
		//fmt.Println(templateStr)
		valueBytes:=rialto.ApplyTemplate(templateStr, deploymentEnvironment)
		var valuesFileName string =  fmt.Sprintf("./values.%v.yaml", plannedServiceInstance.Name)
		err:= ioutil.WriteFile(valuesFileName, valueBytes, 0644)
	    	check(err)//fmt.Println(valueBuf.String())

		f.WriteString(fmt.Sprintf("helm install %v --values %v  \n",
			fmt.Sprintf("%v/%v", plannedServiceInstance.ChartSource.HelmRepo, plannedServiceInstance.ChartSource.ChartName),
			valuesFileName))





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
