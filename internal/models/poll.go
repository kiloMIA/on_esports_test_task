package models

type Poll struct {
    ID       string 
    Question string 
    Options  []string 
    Votes    map[string]int 
}