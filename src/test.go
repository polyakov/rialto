package main

import "rialto"
import "fmt"
import "encoding/json"
import "io/ioutil"

func check(e error) {
	if e != nil {
		panic(e)
	}
}

func main() {
    /*
    var wellKnownTypesFile string = "../test/well-known-types.json"
    rialto.ReadWellKnownServiceConfig( wellKnownTypesFile )
    */

    var services= map[string]rialto.ServiceConfig{}

   /* var service1File string = "../helm/repo/charts/openid-connect-ldap-mitre/www-chart.json"
    service1:= rialto.ReadServiceConfig( service1File )
    bolB, _ := json.Marshal(service1)
    fmt.Println(string(bolB))
    services[service1.Name]=service1
*/
    var service2File string = "../helm/repo/charts/ldap.ApacheDS/ldap-chart.json"
    service2:= rialto.ReadServiceConfig( service2File )
    services[service2.Name]=service2

    var installConfigFile string = "../test/install-env.json"
    installConfig:= rialto.ReadInstallConfig( installConfigFile )
    bolB, _ := json.Marshal( installConfig )
    fmt.Println(string(bolB))


    configFiles := map[string]map[string]string{}

    serviceInstances:=map[string]rialto.ServiceInstance{}
	//values delete
	var val1TemplatePath string = "../helm/repo/charts/ldap.ApacheDS/values.template.yaml"

	// values end delete

    //find properties needed to define
    for _,element := range installConfig.Instances {
        // index is the index where we are
        // element is the element from someSlice for where we are
        //fmt.Println(index)
        //fmt.Println(element.Name)
        foundServiceConfig, ok:= services[element.ServiceName]

        if !ok {
            errStr:=fmt.Sprintf("Service not found: %v", element.Name)
            //fmt.Println(errStr)
            panic(errStr)
        }
        fmt.Println(foundServiceConfig.Name)

        //create service instance from config
        serviceInstances[element.Name]=rialto.GetServiceInstance(foundServiceConfig, element.Name )
	    fmt.Println(serviceInstances[element.Name])
	    var templateStr string = rialto.LoadValueTemplate(element)
	    fmt.Println(templateStr)
	    // load values file
	    var env rialto.Deployment=rialto.Deployment{}
	    env.ServiceInstance=serviceInstances[element.Name]
	    env.Environment=rialto.EnvironmentInstance{}
	    env.Environment.ExternalIPs=[]string{"asdas", "asdsadasdasd"}
	    //rialto.ApplyTemplate(val1TemplatePath,serviceInstances[element.Name])

	    valueBytes:=rialto.ApplyTemplate(val1TemplatePath, env)
	    err:= ioutil.WriteFile(fmt.Sprintf("./values.%v.yaml", element.Name), valueBytes, 0644)
	    check(err)//fmt.Println(valueBuf.String())




        configFile := map[string]string{}
        for _,interfaceVal := range foundServiceConfig.Exposes {
            for _,propVal := range interfaceVal.Properties {
//                fmt.Println(propVal.Name)
                configFile[propVal.Name]="???"
            }
        }

//        fmt.Println(configFile)
        configFiles[foundServiceConfig.Name]=configFile

    }

//    fmt.Println(configFiles)
}
