package main

import (
	"github.com/aws/constructs-go/constructs/v10"
	"github.com/aws/jsii-runtime-go"
	"github.com/cdk8s-team/cdk8s-core-go/cdk8s/v2"
	"github.com/cdk8s-team/cdk8s-plus-go/cdk8splus22/v2"
)

type MyChartProps struct {
	cdk8s.ChartProps
}

type NginxChartProps struct {
	cdk8s.ChartProps
}

var sharedPVC cdk8splus22.PersistentVolumeClaim

// MySQL Chart
func NewMySQLChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	// Creating New Secret
	secretName := "mysqlpassword"
	rootpassword := "VerySecurePassword123"
	user := "kazemi"
	password := "SecurePassword123"
	databasename := "hamravesh"

	mysqlSecret := cdk8splus22.NewSecret(chart, jsii.String("mysql-secret"),
		&cdk8splus22.SecretProps{
			Metadata: &cdk8s.ApiObjectMetadata{Name: jsii.String(secretName)}})

	secretKey := "VerySecureSecretKey123"

	mysqlSecret.AddStringData(jsii.String(secretKey), jsii.String(rootpassword))
	mysqlSecret.AddStringData(jsii.String(secretKey), jsii.String(user))
	mysqlSecret.AddStringData(jsii.String(secretKey), jsii.String(password))
	mysqlSecret.AddStringData(jsii.String(secretKey), jsii.String(databasename))

	deploymentName := "nextcloud"

	// MySQL Deployment
	mysqldep := cdk8splus22.NewDeployment(chart, jsii.String("mysql-deployment"), &cdk8splus22.DeploymentProps{Metadata: &cdk8s.ApiObjectMetadata{Labels: &map[string]*string{"app": &deploymentName}}})

	containerImage := "mysql"

	mysqlContainer := mysqldep.AddContainer(&cdk8splus22.ContainerProps{
		Name:  jsii.String("mysql-container"),
		Image: jsii.String(containerImage),
		// Port:  jsii.Number(3306),
	})

	// Using Secret in Env
	envValRootPassFromSecret := cdk8splus22.EnvValue_FromSecretValue(&cdk8splus22.SecretValue{Key: jsii.String(secretKey), Secret: mysqlSecret}, &cdk8splus22.EnvValueFromSecretOptions{Optional: jsii.Bool(true)})
	envValUserFromSecret := cdk8splus22.EnvValue_FromSecretValue(&cdk8splus22.SecretValue{Key: jsii.String(secretKey), Secret: mysqlSecret}, &cdk8splus22.EnvValueFromSecretOptions{Optional: jsii.Bool(true)})
	envValPassFromSecret := cdk8splus22.EnvValue_FromSecretValue(&cdk8splus22.SecretValue{Key: jsii.String(secretKey), Secret: mysqlSecret}, &cdk8splus22.EnvValueFromSecretOptions{Optional: jsii.Bool(true)})
	envValDBFromSecret := cdk8splus22.EnvValue_FromSecretValue(&cdk8splus22.SecretValue{Key: jsii.String(secretKey), Secret: mysqlSecret}, &cdk8splus22.EnvValueFromSecretOptions{Optional: jsii.Bool(true)})

	mySQLUserEnvName := "MYSQL_USER"
	mySQLPasswordEnvName := "MYSQL_PASSWORD"
	mySQLRootPasswordEnvName := "MYSQL_ROOT_PASSWORD"
	mySQLDBEnvName := "MYSQL_DATABASE"

	mysqlContainer.Env().AddVariable(jsii.String(mySQLRootPasswordEnvName), envValRootPassFromSecret)
	mysqlContainer.Env().AddVariable(jsii.String(mySQLUserEnvName), envValUserFromSecret)
	mysqlContainer.Env().AddVariable(jsii.String(mySQLPasswordEnvName), envValPassFromSecret)
	mysqlContainer.Env().AddVariable(jsii.String(mySQLDBEnvName), envValDBFromSecret)

	// Using Shared PVC
	mysqlVolumeName := "mysql-persistent-storage"
	mysqlVolume := cdk8splus22.Volume_FromPersistentVolumeClaim(chart, jsii.String("mysql-vol-pvc"), sharedPVC, &cdk8splus22.PersistentVolumeClaimVolumeOptions{Name: jsii.String(mysqlVolumeName)})

	mysqlVolumeMountPath := "/var/lib/mysql"
	mysqlVolumeMountSubPath := "server-data"

	mysqlContainer.Mount(jsii.String(mysqlVolumeMountPath), mysqlVolume, &cdk8splus22.MountOptions{SubPath: &mysqlVolumeMountSubPath})

	// Creating Service for MySQL Deployment
	mysqlServiceName := "mysql-service"
	// clusterIPNone := "None"
	appName := "nextcloud"

	cdk8splus22.NewService(chart, jsii.String("mysql-service"), &cdk8splus22.ServiceProps{
		Metadata: &cdk8s.ApiObjectMetadata{Name: jsii.String(mysqlServiceName), Labels: &map[string]*string{"app": &appName}},
		Selector: mysqldep,
		// ClusterIP: jsii.String(clusterIPNone),
		Ports: &[]*cdk8splus22.ServicePort{{Port: jsii.Number(3306)}},
	})

	return chart
}

// Nextcloud Chart
func NewNextCloudChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	deploymentName := "nextcloud"
	// NextCloud Deployment
	nextclouddep := cdk8splus22.NewDeployment(chart, jsii.String("nextcloud-deployment"), &cdk8splus22.DeploymentProps{Metadata: &cdk8s.ApiObjectMetadata{Labels: &map[string]*string{"app": &deploymentName}}})

	containerImage := "nextcloud:16-apache"

	nextcloudContainer := nextclouddep.AddContainer(&cdk8splus22.ContainerProps{
		Name:  jsii.String("nextcloud-container"),
		Image: jsii.String(containerImage),
		// Port:  jsii.Number(80),
	})

	// Using Shared PVC
	nextcloudVolumeName := "nextcloud-persistent-storage"
	nextcloudVolume := cdk8splus22.Volume_FromPersistentVolumeClaim(chart, jsii.String("nextcloud-vol-pvc"), sharedPVC, &cdk8splus22.PersistentVolumeClaimVolumeOptions{Name: jsii.String(nextcloudVolumeName)})

	nextcloudVolumeMountPath := "/var/www/html"
	nextcloudVolumeMountSubPath := "server-data"

	nextcloudContainer.Mount(jsii.String(nextcloudVolumeMountPath), nextcloudVolume, &cdk8splus22.MountOptions{SubPath: &nextcloudVolumeMountSubPath})

	nextcloudServiceName := "nextcloud-service"
	// clusterIPNone := "None"
	appName := "nextcloud"

	cdk8splus22.NewService(chart, jsii.String("nextcloud-service"), &cdk8splus22.ServiceProps{
		Metadata: &cdk8s.ApiObjectMetadata{Name: jsii.String(nextcloudServiceName), Labels: &map[string]*string{"app": &appName}},
		Selector: nextclouddep,
		// ClusterIP: jsii.String(clusterIPNone),
		Ports: &[]*cdk8splus22.ServicePort{{Port: jsii.Number(80)}},
	})

	return chart
}

// LoadBalancer Chart
func NewLoadBalancerChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	loadbalancerServiceName := "nextcloud-service"
	appName := "nextcloud"

	cdk8splus22.NewService(chart, jsii.String("nextcloud-service"), &cdk8splus22.ServiceProps{
		Metadata: &cdk8s.ApiObjectMetadata{Name: jsii.String(loadbalancerServiceName), Labels: &map[string]*string{"app": &appName}},
		// Selector: NewNextCloudChart(),
		Type: cdk8splus22.ServiceType_LOAD_BALANCER,
		// ClusterIP: jsii.String(clusterIPNone),
		Ports: &[]*cdk8splus22.ServicePort{{Port: jsii.Number(80)}},
	})

	return chart
}

// PVC Chart
func NewPVCChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	pvcName := "shared-pvc"
	appName := "nextcloud"
	storageClassType := "rawfile-btrfs"

	sharedPVC = cdk8splus22.NewPersistentVolumeClaim(chart, jsii.String("shared-pvc"), &cdk8splus22.PersistentVolumeClaimProps{
		Metadata:         &cdk8s.ApiObjectMetadata{Name: jsii.String(pvcName), Labels: &map[string]*string{"app": &appName}},
		AccessModes:      &[]cdk8splus22.PersistentVolumeAccessMode{cdk8splus22.PersistentVolumeAccessMode_READ_WRITE_ONCE},
		Storage:          cdk8s.Size_Gibibytes(jsii.Number(1)),
		StorageClassName: &storageClassType,
		VolumeMode:       cdk8splus22.PersistentVolumeMode_FILE_SYSTEM,
	})

	return chart
}

// Nginx Chart
func NewNginxChart(scope constructs.Construct, id string, props *NginxChartProps) cdk8s.Chart {
	var cprops cdk8s.ChartProps
	if props != nil {
		cprops = props.ChartProps
	}
	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

	nginxdep := cdk8splus22.NewDeployment(chart, jsii.String("deployment"), &cdk8splus22.DeploymentProps{Metadata: &cdk8s.ApiObjectMetadata{Name: jsii.String("nginx-deployment")}})

	nginxdep.AddContainer(&cdk8splus22.ContainerProps{
		Name:  jsii.String("nginx-container"),
		Image: jsii.String("nginx"),
		// Port:  jsii.Number(80),
	})

	nginxdep.ExposeViaService(&cdk8splus22.DeploymentExposeViaServiceOptions{
		Name:        jsii.String("nginx-container-service"),
		ServiceType: cdk8splus22.ServiceType_LOAD_BALANCER,
		Ports:       &[]*cdk8splus22.ServicePort{{Port: jsii.Number(9090), TargetPort: jsii.Number(80)}}})

	return chart
}

// // Ingress Chart
// func NewIngressChart(scope constructs.Construct, id string, props *MyChartProps) cdk8s.Chart {
// 	var cprops cdk8s.ChartProps
// 	if props != nil {
// 		cprops = props.ChartProps
// 	}
// 	chart := cdk8s.NewChart(scope, jsii.String(id), &cprops)

// 	ingressName := "cluster-ingress"
// 	host := "nextcloud.site.test"
// 	// NextCloud Deployment
// 	nginxingress := cdk8splus22.NewIngress(chart, jsii.String("nextcloud-deployment"), &cdk8splus22.IngressProps{Metadata: &cdk8s.ApiObjectMetadata{Labels: &map[string]*string{"app": &ingressName}}})
// 	nginxingress.AddRules(&cdk8splus22.IngressRule{Host: &host})
// 	nginxingress.AddDefaultBackend(cdk8splus22.IngressBackend_FromService)
// 	return chart
// }

func main() {
	app := cdk8s.NewApp(nil)

	// ingressChart := NewIngressChart(app, "ingress", nil)
	loadbalancerChart := NewLoadBalancerChart(app, "loadbalancer", nil)
	sharedPVCChart := NewPVCChart(app, "pvc", nil)
	mysqlChart := NewMySQLChart(app, "mysql", nil)
	nextcloudChart := NewNextCloudChart(app, "nextcloud", nil)
	nginxChart := NewNginxChart(app, "nginx", nil)
	mysqlChart.AddDependency(sharedPVCChart, nginxChart)
	nextcloudChart.AddDependency(sharedPVCChart, mysqlChart, nginxChart, loadbalancerChart)
	app.Synth()
}
