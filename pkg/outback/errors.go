package outback

var (
	errFailedToListClusters     = "error listing clusters"
	errFailedToListServices     = "error listing services"
	errFailedToListRunningTasks = "error listing running tasks"

	errCouldNotRetrieveCluster        = "could not retrieve cluster"
	errCouldNotRetrieveService        = "could not retrieve service"
	errCouldNotRetrieveTaskDefinition = "could not retrieve task definition"
	errCouldNotRetrieveTasks          = "could not retrieve tasks"
	errCouldNotRetrieveImages         = "could not retrieve images"

	errInvalidTaskDefinition = "task definition contains no container definitions"

	errCouldNotRegisterTaskDefinition = "could not register new task definition"
	errCouldNotUpdateService          = "could not update service"

	errClusterNotFound = "cluster was not found"
	errServiceNotFound = "service was not found"

	errCouldNotRunTask = "desired task could not run"

	errCouldNotGetLogs = "could not get cloudwatch logs"

	errECRLogin = "Could not login to ECR"
)
