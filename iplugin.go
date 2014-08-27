package main

type IPlugin interface {
    Init() error
}