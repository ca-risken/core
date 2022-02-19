package main

type projectService struct {
	repository projectRepository
	iamClient  iamService
}
