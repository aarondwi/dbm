package main

import (
	"testing"
)

var dbMock = &dummyDB{}
var dbMockFail = &dummyDBFail{}
var sourceMock = &dummySource{}
var sourceMockFail = &dummySourceFail{}
var sourceMockFailOnReadSrc = &dummySourceFailOnReadSrc{}

func TestInit(t *testing.T) {
	dirname := "test"
	err := Init(sourceMock, dirname)
	if err != nil {
		t.Fatalf(`Init function failed: %v`, err)
	}

	err = Init(sourceMockFail, dirname)
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

	err = GenerateSrcfile(sourceMockFail, filename)
	if err == nil {
		t.Fatal(`GenerateSrcfile function should fail but it is not`)
	}
}

func TestReadConfigFile(t *testing.T) {
	result, err := ReadConfigFile(sourceMock)
	if err != nil {
		t.Fatalf("ReadConfigFile failed: %v,", err)
	}
	if result.Dialect != "test" ||
		result.Host != "test" ||
		result.Port != 1000 ||
		result.Username != "test" ||
		result.Password != "test" ||
		result.Database != "test" ||
		result.Sslmode != "test" {
		t.Fatalf(
			"The Content is different: %v", result)
	}

	result, err = ReadConfigFile(sourceMockFail)
	if err == nil {
		t.Fatal(`TestReadConfigFile with mockFail should fail but it is not`)
	}
}

func TestSetup(t *testing.T) {
	err := Setup(dbMock)
	if err != nil {
		t.Fatalf("CreateLogTable failed: %v", err)
	}

	err = Setup(dbMockFail)
	if err == nil {
		t.Fatal(`TestSetup with mockFail should fail but it is not`)
	}
}

func TestStatus(t *testing.T) {
	err := Status(sourceMock, dbMock)
	if err != nil {
		t.Fatalf("Status failed: %v", err)
	}

	err = Status(sourceMockFail, dbMock)
	if err == nil {
		t.Fatal(`TestSetupFail function should fail but it is not`)
	}

	err = Status(sourceMock, dbMockFail)
	if err == nil {
		t.Fatal(`TestSetupFail function should fail but it is not`)
	}
}

func TestUp(t *testing.T) {
	err := Up(sourceMock, dbMock, "")
	if err != nil {
		t.Fatalf("TestUp with empty filename fail: %v", err)
	}

	err = Up(sourceMock, dbMock, "somefile")
	if err == nil {
		t.Fatal("TestUp should fail because already applied, but it is not")
	}

	err = Up(sourceMock, dbMock, "anotherfile")
	if err != nil {
		t.Fatalf("TestUp with literal 'filename' fail: %v", err)
	}

	err = Up(sourceMock, dbMock, "notfound")
	if err == nil {
		t.Fatal("TestUp should fail because 'notfound' but it is not")
	}

	err = Up(sourceMockFail, dbMock, "")
	if err == nil {
		t.Fatal("TestUp should fail because can't read from src dir, but it is not")
	}

	err = Up(sourceMock, dbMockFail, "")
	if err == nil {
		t.Fatal("TestUp should fail because can't read log from db, but it is not")
	}

	err = Up(sourceMockFailOnReadSrc, dbMock, "")
	if err == nil {
		t.Fatal("TestUp should fail because can't read log from db, but it is not")
	}

	err = Up(sourceMockFailOnReadSrc, dbMock, "anotherfile")
	if err == nil {
		t.Fatal("TestUp should fail because can't read log from db, but it is not")
	}
}

func TestDown(t *testing.T) {
	err := Down(sourceMock, dbMock, "")
	if err != nil {
		t.Fatalf("TestDown with empty filename fail: %v", err)
	}

	err = Down(sourceMock, dbMock, "somefile")
	if err != nil {
		t.Fatalf("TestDown with filename 'somefail' fail: %v", err)
	}

	err = Down(sourceMock, dbMock, "anotherfile")
	if err == nil {
		t.Fatalf("TestDown with filename 'anotherfile' should fail because it has not been applied but it is not")
	}

	err = Down(sourceMock, dbMock, "notfound")
	if err == nil {
		t.Fatalf("TestDown with filename 'notfound' should fail because it is not in src dir, but it is not")
	}

	err = Down(sourceMock, dbMockFail, "")
	if err == nil {
		t.Fatalf("TestDown should fail because can't read from db, but it is not")
	}

	err = Down(sourceMockFail, dbMock, "")
	if err == nil {
		t.Fatalf("TestDown with empty filename should fail because can't read file, but it is not")
	}

	err = Down(sourceMockFail, dbMock, "somefile")
	if err == nil {
		t.Fatalf("TestDown with filename 'somefile' should fail because can't read file, but it is not")
	}
}
