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

func prompt(message string, defaultVal string)string {
	fmt.Printf(message)
	var input string=defaultVal
	fmt.Scanln(&input)

	return input;

}
var context Context

func Init() Context {

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

	return loadUrl(url)

}

func LoadServiceConfig(serviceInstance ServiceInstance) ServiceConfig {
	var url string = serviceInstance.ChartSource.Url+
		"/"+serviceInstance.ChartSource.ChartName+
		"/hspc.json"
	//fmt.Println(url)
	url=strings.TrimSpace( url )

	strConfig:=loadUrl(url)

	return ParseServiceConfig(strConfig)
}

func GetExternalIPs()[]string{
	var ips []string = []string{}
	for nextLoopOk := true; nextLoopOk;
	{
		var newIP string = prompt("Enter external IP. (Enter to continue)","")
		if(newIP== "") {
			nextLoopOk=false;
		} else {
			ips = append(ips, newIP)
		}
	}

	return ips
}

func loadUrl(url string)string{
	fmt.Println(url)
	res, err := http.Get(url)
	if(res.StatusCode != 200) {
		panic(fmt.Sprintf("Request to %v, response code: %v\n", url, res.StatusCode))
	}
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
			input:=prompt(fmt.Sprintf("Exposes: Enter value for %v.%v[%v]\n", foundServiceConfig.Name, foundProperty.Name, foundProperty.DefaultValue), foundProperty.DefaultValue)
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
			input:=prompt(fmt.Sprintf("DependsOn: Enter value for %v.%v[%v]\n", foundInterface.Name, foundProperty.Name, foundProperty.DefaultValue), foundProperty.DefaultValue)
			dependsOnProperties.Properties[foundProperty.Name]=input;
		}
		serviceInstance.DependsOn[foundInterface.Name] = dependsOnProperties
	}

	var serviceProperties=map[string]string{}
	for _,foundProperty:= range foundServiceConfig.Properties {
		input:=prompt(fmt.Sprintf("Service Properties: Enter value for %v.%v[%v]\n", foundServiceConfig.Name, foundProperty.Name, foundProperty.DefaultValue), foundProperty.DefaultValue)
		serviceProperties[foundProperty.Name]=input;
	}
	serviceInstance.Properties = serviceProperties
	return serviceInstance;
}

//ReadServiceConfig - parse JSON service config and return ServiceConfig stuct
func ReadServiceConfig( filepath string)  ServiceConfig {
	dat, err := ioutil.ReadFile(filepath)
	check(err)
	var jsonStr = string(dat)
//      fmt.Print(jsonStr)
	return ParseServiceConfig(jsonStr);
}

func ParseServiceConfig( content string)  ServiceConfig {
	res := ServiceConfig{}
	json.Unmarshal([]byte(content), &res)
	//fmt.Println(res)
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
