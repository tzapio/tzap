package project

type ProjectDir string
type ProjectName string

type ProjectDB map[ProjectName]ProjectDir

const LOCALPROJECTNAME ProjectName = "@LOCAL"
