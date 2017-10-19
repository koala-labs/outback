package main

//import (
//	"fmt"
//
//	"github.com/aws/aws-sdk-go/aws/awserr"
//	"github.com/aws/aws-sdk-go/service/ecr"
//	"github.com/aws/aws-sdk-go/service/ecs"
//)
//
//func handleECSErr(err error) {
//	if err != nil {
//		if aerr, ok := err.(awserr.Error); ok {
//			switch aerr.Code() {
//			case ecs.ErrCodeServerException:
//				fmt.Println(ecs.ErrCodeServerException, aerr.Error())
//			case ecs.ErrCodeClientException:
//				fmt.Println(ecs.ErrCodeClientException, aerr.Error())
//			case ecs.ErrCodeInvalidParameterException:
//				fmt.Println(ecs.ErrCodeInvalidParameterException, aerr.Error())
//			case ecs.ErrCodeClusterNotFoundException:
//				fmt.Println(ecs.ErrCodeClusterNotFoundException, aerr.Error())
//			case ecs.ErrCodeServiceNotFoundException:
//				fmt.Println(ecs.ErrCodeServiceNotFoundException, aerr.Error())
//			case ecs.ErrCodeServiceNotActiveException:
//				fmt.Println(ecs.ErrCodeServiceNotActiveException, aerr.Error())
//			default:
//				fmt.Println(aerr.Error())
//			}
//		} else {
//			fmt.Println(err.Error())
//		}
//		return
//	}
//}
//
//func handleECRErr(err error) {
//	if err != nil {
//		if aerr, ok := err.(awserr.Error); ok {
//			switch aerr.Code() {
//			case ecr.ErrCodeServerException:
//				fmt.Println(ecr.ErrCodeServerException, aerr.Error())
//			case ecr.ErrCodeInvalidParameterException:
//				fmt.Println(ecr.ErrCodeInvalidParameterException, aerr.Error())
//			case ecr.ErrCodeRepositoryNotFoundException:
//				fmt.Println(ecr.ErrCodeRepositoryNotFoundException, aerr.Error())
//			default:
//				fmt.Println(aerr.Error())
//			}
//		} else {
//			fmt.Println(err.Error())
//		}
//		return
//	}
//}
