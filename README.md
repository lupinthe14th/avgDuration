# Average Task Duration Calculator

A Golang tool that parses nested JSON input representing tasks, computes, and outputs the average duration between task start and creation times. Includes modular functions for easy testing.

## Overview:

The program processes a nested JSON structure that contains information about tasks, including their start and creation times. The main goal is to compute the average duration between the start and creation times of these tasks.

1. **Data Structure**:

    - A `Task` structure represents individual task data, with fields for `StartedAt`, `CreatedAt`, and `TaskArn`.

1. **Functionality**:
    - `AverageDuration`: This function accepts a list of tasks and computes the average duration between the start and creation times. It parses the timestamps, calculates the difference between the start and creation times for each task, sums up these differences, and finally calculates the average duration.
1. **Main Execution**:

    - The main function reads the JSON data from standard input.
    - It then unmarshals this data into the nested task structure.
    - The average duration is calculated using the `AverageDuration` function.
    - The program prints the calculated average duration.

1. **Testing**:
    - A separate test file (`avgDuration_test.go`) provides unit tests for the `AverageDuration` function to ensure its accuracy.
