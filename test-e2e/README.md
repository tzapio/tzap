# Test-e2e Directory

The **test-e2e** directory contains end-to-end tests for our application. These tests simulate real user scenarios and help ensure that all the components of our application are working together correctly.

The purpose of the test-e2e directory is to provide a comprehensive set of tests that cover the entire workflow of our application. These tests help us identify any issues or bugs that may arise due to interactions between different parts of our system.

In the test-e2e directory, you will find a collection of test files that are organized into different categories based on the functionality they test. Each test file contains a set of test cases that verify specific features or use cases of our application.

To run the end-to-end tests, you can use the test runner provided in the test-e2e directory. This test runner executes all the test files and reports any failures or errors encountered during the testing process.

We use Docker to build an image of our application for testing. The Dockerfile is located in the root of the test-e2e directory. The Docker image is built using the command `docker build --build-arg OPENAI_APIKEY=${OPENAI_APIKEY} . -t tzap-test-e2e` as specified in the `.github/workflows/test-e2e.yml` file.

The tests are executed in the Docker container using the command `docker run --rm tzap-test-e2e make test-refactor` as specified in the `.github/workflows/test-e2e.yml` file.

We recommend regularly running the end-to-end tests as part of your development workflow to ensure the overall integrity and functionality of our application.

Please refer to the documentation in the test-e2e directory for more information on how to run and analyze the results of the end-to-end tests.