package rialto

import "fmt"
import "encoding/json"
//import "os"
import "bytes"
//import "bufio"
import "io/ioutil"
import "net/http"
import "text/template"
//import "github.com/Masterminds/sprig"
import "strings"

func check(e error) {
    if e != nil {
        panic(e)
    }
}

func prompt(message string)string {
	fmt.Printf(message)
	var input string
	fmt.Scanln(&input)

	return input;

}
var context Context

func InitMe() Context {

	context.Services=map[string]ServiceConfig{}
	// load service catalog
	var service1File string = "../helm/repo/charts/openid-connect-ldap-mitre/www-chart.json"
	service1:= ReadServiceConfig( service1File )
 //bolB, _ := json.Marshal(service1)
 //fmt.Println(string(bolB))
 	context.Services[service1.Name]=service1

	var service2File string = "../helm/repo/charts/ldap.ApacheDS/ldap-chart.json"
	service2:= ReadServiceConfig( service2File )
	context.Services[service2.Name]=service2

	return context
}

func LoadValueTemplate(serviceInstance ServiceInstance) string {
	var url string = serviceInstance.ChartSource.Url+
		"/"+serviceInstance.ChartSource.ChartName+
		"/values.template.yaml"
	//fmt.Println(url)
	url=strings.TrimSpace( url )

	res, err := http.Get(url)
	check(err)

	valueTemplate, err := ioutil.ReadAll(res.Body)
	defer res.Body.Close()
	check(err)

	//fmt.Printf("valueTemplate: %s", valueTemplate)
	str := string(valueTemplate[:])
	return str

}

// eventully we want to be able to say "I need a FHIR server" and have the system find as service to fit this requirement
/*
func ReadWellKnownServiceConfig( filepath string ) WellKnownServiceConfig {
    fmt.Printf("filename: %v", filepath)

    dat, err := ioutil.ReadFile(filepath)
    check(err)

    //var jsonStr = string(dat)
    //fmt.Print(jsonStr)

    res := WellKnownServiceConfig{}
    json.Unmarshal(dat, &res)
    fmt.Println(res)

    return res
}
*/

func ApplyTemplate( valStr string, deployment Deployment)[]byte {
	/*
	dat, err := ioutil.ReadFile(filepath)
	check(err)
	var valStr string = string(dat)
	*/

	tmpl, err := template.New("value").Parse(valStr)
	check(err)

	var outBuf bytes.Buffer;
	err = tmpl.Execute(&outBuf, deployment)
	check(err)
	//fmt.Println("==============")
	//fmt.Println(outBuf)
	return outBuf.Bytes()
}


func GetServiceInstance(foundServiceConfig ServiceConfig, name string) ServiceInstance {
	var serviceInstance ServiceInstance=ServiceInstance{}
	serviceInstance.Name=name
	serviceInstance.ServiceName=foundServiceConfig.Name
	serviceInstance.Exposes=map[string]ServiceInstanceProperties{}
	serviceInstance.DependsOn=map[string]ServiceInstanceProperties{}
	serviceInstance.Properties=map[string]string{}


	// map Exposes
	for _,foundInterface:= range foundServiceConfig.Exposes {
		//fmt.Printf("Found exposes interface: %v: %v\n", interfaceIndex, foundInterface.Name)
		var exposedProperties ServiceInstanceProperties = ServiceInstanceProperties{}
		exposedProperties.Properties = map[string]string{}
		for _, foundProperty := range foundInterface.Properties {
			input:=prompt(fmt.Sprintf("Enter value for %v.%v[Exposes]\n", foundServiceConfig.Name, foundProperty.Name))
			exposedProperties.Properties[foundProperty.Name] = input;
		}
		serviceInstance.Exposes[foundInterface.Name] = exposedProperties
	}

	for _,foundInterface:= range foundServiceConfig.DependsOn {
	//	fmt.Printf("Found dependsOn: %v: %v\n", interfaceIndex, foundInterface.Name)
		var dependsOnProperties ServiceInstanceProperties = ServiceInstanceProperties{}
		dependsOnProperties.Properties=map[string]string{}
		for _,foundProperty:= range foundInterface.Properties {
	//		fmt.Printf("Found prop: %v: %v\n", indexProp, foundProperty.Name)
			input:=prompt(fmt.Sprintf("Enter value for %v.%v[DependsOn]\n", foundInterface.Name, foundProperty.Name))
			dependsOnProperties.Properties[foundProperty.Name]=input;
		}
		serviceInstance.DependsOn[foundInterface.Name] = dependsOnProperties
	}

	/*for interfaceIndex,foundInterface:= range foundServiceConfig.Properties {
		fmt.Printf("Found properties")
		var properties ServiceInstanceProperties = ServiceInstanceProperties{}
	*/
	var serviceProperties=map[string]string{}
	for _,foundProperty:= range foundServiceConfig.Properties {
	//	fmt.Printf("Found prop: %v: %v\n", indexProp, foundProperty.Name)
		input:=prompt(fmt.Sprintf("Enter value for service property %v.%v\n", foundServiceConfig.Name, foundProperty.Name))
		serviceProperties[foundProperty.Name]=input;
	}
	serviceInstance.Properties = serviceProperties

	return serviceInstance;
}

//ReadServiceConfig - parse JSON service config and return ServiceConfig stuct
func ReadServiceConfig( filepath string)  ServiceConfig {
    //fmt.Printf("filename: %v", filepath)

    dat, err := ioutil.ReadFile(filepath)
    check(err)

//    var jsonStr = string(dat)
//    fmt.Print(jsonStr)

    res := ServiceConfig{}
    json.Unmarshal(dat, &res)
    fmt.Println(res)

    return res;
}

//ReadInstallConfig - read a JSON doc with configuration of service instances
func ReadInstallConfig( filepath string)  InstallConfig {
    //fmt.Printf("filename: %v", filepath)

    dat, err := ioutil.ReadFile(filepath)
    check(err)

    //var jsonStr = string(dat)
    //fmt.Print(jsonStr)

    res := InstallConfig{}
    json.Unmarshal(dat, &res)
//    fmt.Println(res)

    return res;
}
