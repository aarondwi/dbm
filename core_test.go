package main

import (
	"testing"
)

var dbMock = &dummyDB{}
var dbMockFail = &dummyDBFail{}
var sourceMock = &dummySource{}
var sourceMockFail = &dummySourceFail{}

func TestInit(t *testing.T) {
	dirname := "test"
	err := Init(sourceMock, dirname)
	if err != nil {
		t.Fatalf(`Init function failed: %v`, err)
	}
}

func TestInitFail(t *testing.T) {
	dirname := "test"
	err := Init(sourceMockFail, dirname)
	if err == nil {
		t.Fatal(`Init function should fail but it is not`)
	}
}

func TestGenerateSrcfile(t *testing.T) {
	filename := "CreateTableDummy"
	err := GenerateSrcfile(sourceMock, filename)
	if err != nil {
		t.Fatalf("GenerateSrcfile failed: %v", err)
	}
}

func TestGenerateSrcfileFail(t *testing.T) {
	filename := "CreateTableDummy"
	err := GenerateSrcfile(sourceMockFail, filename)
	if err == nil {
		t.Fatal(`GenerateSrcfile function should fail but it is not`)
	}
}

func TestSetup(t *testing.T) {
	err := Setup(dbMock)
	if err != nil {
		t.Fatalf("CreateLogTable failed: %v", err)
	}
}

func TestSetupFail(t *testing.T) {
	err := Setup(dbMockFail)
	if err == nil {
		t.Fatal(`TestSetupFail function should fail but it is not`)
	}
}

func TestStatus(t *testing.T) {
	err := Status(sourceMock, dbMock)
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}
}

func TestStatusReadDirFail(t *testing.T) {
	err := Status(sourceMockFail, dbMock)
	if err == nil {
		t.Fatal(`TestSetupFail function should fail but it is not`)
	}
}

func TestStatusReadDBFail(t *testing.T) {
	err := Status(sourceMock, dbMockFail)
	if err == nil {
		t.Fatal(`TestSetupFail function should fail but it is not`)
	}
}

func TestUp(t *testing.T) {

}

func TestDown(t *testing.T) {

}
