package ufo

import "errors"

var (
	ErrFailedToListClusters = errors.New("Failed to list clusters.")
	ErrFailedToListServices = errors.New("Failed to list services.")
	ErrFailedToListRunningTasks = errors.New("Failed to list running tasks.")

	ErrCouldNotRetrieveCluster = errors.New("Could not retrieve cluster.")
	ErrCouldNotRetrieveService = errors.New("Could not retrieve service.")
	ErrCouldNotRetrieveTaskDefinition = errors.New("Could not retrieve task definition.")
	ErrCouldNotRetrieveTasks = errors.New("Could not retrieve tasks.")
	ErrCouldNotRetrieveImages = errors.New("Could not retrieve images.")

	ErrInvalidTaskDefinition = errors.New("Invalid task definition for operation.")

	ErrCouldNotRegisterTaskDefinition = errors.New("Could not register new task definition.")
	ErrCouldNotUpdateService = errors.New("Could not update service.")

	ErrClusterNotFound = errors.New("Cluster was not found.")
	ErrServiceNotFound = errors.New("Service was not found.")
)