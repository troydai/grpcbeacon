# CI/CD Setup with Integration Tests

This document explains the comprehensive CI/CD pipeline setup for the gRPC Beacon service, including how integration tests are integrated into the GitHub Actions workflow.

## 🚀 **Overview**

The CI/CD pipeline ensures code quality, runs comprehensive tests (including integration tests), and builds secure Docker images. The pipeline is designed to fail fast if any tests fail, preventing broken code from being deployed.

## 📋 **Workflow Structure**

### 1. **Main CI Workflow** (`ci.yml`)

**Triggers:**
- ✅ Every push to `main` branch
- ✅ Every pull request to `main` branch

**Jobs:**

#### **Test Job**
- **Go Setup**: Uses Go 1.24
- **Dependency Caching**: Caches Go modules for faster builds
- **Code Quality**: Runs `go vet` and `go fmt` checks
- **Unit Tests**: Runs all unit tests with race detection and coverage
- **🎯 Integration Tests**: Runs our comprehensive integration tests
- **Coverage Upload**: Uploads coverage reports to Codecov

#### **Build Job** 
- **Depends on**: Test job passing ✅
- **Binary Build**: Builds the `beacon-server` binary
- **Build Verification**: Verifies the binary was created successfully
- **Package Build**: Builds all Go packages

#### **Lint Job**
- **Static Analysis**: Runs `golangci-lint` for comprehensive code analysis
- **Best Practices**: Enforces Go best practices and style guides

#### **Security Job**
- **Security Scanning**: Runs `gosec` for security vulnerability detection
- **SARIF Upload**: Uploads security findings to GitHub Security tab

### 2. **Docker Publishing Workflows**

#### **Server Docker Workflow** (`docker-publish-server.yml`)
- **Pre-Build Testing**: Runs full test suite **INCLUDING INTEGRATION TESTS** before building Docker image
- **Multi-Platform**: Builds for `linux/amd64` and `linux/arm64`
- **Registry**: Publishes to GitHub Container Registry (`ghcr.io`)
- **Security**: Signs images with cosign for supply chain security

#### **Toolbox Docker Workflow** (`docker-publish-toolbox.yml`)
- **Same Testing**: Runs full test suite **INCLUDING INTEGRATION TESTS** before building
- **Multi-Platform**: Builds for `linux/amd64` and `linux/arm64`
- **Registry**: Publishes to GitHub Container Registry (`ghcr.io`)

## 🧪 **Integration Tests in Action**

### **What Gets Tested**
1. **Server Startup**: Verifies the gRPC server starts successfully
2. **Echo Endpoint**: Tests the `Signal` method (echo functionality)
3. **Health Checks**: Verifies gRPC health check endpoints
4. **Concurrent Calls**: Tests multiple simultaneous requests
5. **Error Handling**: Verifies proper error responses
6. **Resource Cleanup**: Ensures proper server shutdown

### **Test Execution Flow**
```bash
# Integration tests run in multiple places:
1. CI Workflow → Test Job → "Run integration tests"
2. Docker Server Workflow → Build Job → "Run tests (including integration tests)"
3. Docker Toolbox Workflow → Build Job → "Run tests (including integration tests)"
```

### **Test Commands**
```bash
# Standard test run (includes integration tests)
go test -v ./...

# Specific integration test run
go test -v -run TestIntegration ./...

# With race detection and coverage
go test -v -race -coverprofile=coverage.out -covermode=atomic ./...
```

## 🔒 **Security & Quality Gates**

### **All workflows include:**
- ✅ **Integration Tests** - Must pass before any Docker build
- ✅ **Unit Tests** - Full test suite with race detection
- ✅ **Code Quality** - Static analysis and formatting checks
- ✅ **Security Scanning** - Vulnerability detection
- ✅ **Build Verification** - Ensures binaries compile correctly
- ✅ **Dependency Verification** - Validates go.mod integrity

### **Fail-Fast Approach**
- If integration tests fail → Docker build is **BLOCKED**
- If any test fails → No images are published
- If security issues found → Reported to GitHub Security tab
- If formatting issues → CI fails with clear error messages

## 📊 **Monitoring & Reporting**

### **Coverage Reports**
- Uploaded to Codecov for trend analysis
- Includes integration test coverage
- Tracks coverage changes over time

### **Security Reports**
- SARIF format uploaded to GitHub Security tab
- Automatic security advisory integration
- Dependabot alerts for vulnerable dependencies

### **Build Artifacts**
- Server binary verification
- Multi-platform Docker images
- Signed container images for supply chain security

## 🎯 **Key Benefits**

1. **Comprehensive Testing**: Integration tests run on every code change
2. **Early Detection**: Issues caught before Docker build/publish
3. **Multi-Platform**: Tests run on GitHub's Ubuntu runners
4. **Caching**: Fast builds with Go module caching
5. **Security**: Signed images and vulnerability scanning
6. **Reliability**: Real gRPC server testing with actual client connections

## 🔧 **Local Testing**

Developers can run the same tests locally:
```bash
# Run all tests (including integration tests)
go test -v ./...

# Run only integration tests
go test -v -run TestIntegration ./...

# Run with race detection
go test -v -race ./...
```

## 🚀 **Deployment Flow**

1. **Code Push/PR** → CI workflow runs
2. **All Tests Pass** → Docker workflows trigger
3. **Integration Tests Pass** → Docker images built
4. **Images Published** → Available for deployment
5. **Security Scan Complete** → Ready for production

This setup ensures that **no code reaches production without passing comprehensive integration tests**, providing confidence in the stability and reliability of the gRPC Beacon service.