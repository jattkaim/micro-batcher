# Minimal Batching library
The purpose of this library is to quickly allow applications to batch requests that are being run at a high throughput.

# Getting started

### Prerequisites
- Go installed on your system (version 1.11 or higher recommended for module support).

### Installation
```
go get github.com/jattkaim/micro-batcher
```

### Usage

1. **Implement the Processor Interface**: This library requires a `Processor` interface to be implemented by the user, responsible for processing jobs. Define this processor in your application before using the micro-batcher library:

    ```go
    type MyAwesomeBatchProcessor struct {
        // Your batch processor fields here
    }

    func (p *MyAwesomeBatchProcessor) Process(jobs []batcher.Job) ([]batcher.JobResult, error) {
        // Implement your batch processing logic here
    }
	
    func (p *MyAwesomeBatchProcessor) Start() error {
        return nil // return nil if not needed
    }

    func (p *MyAwesomeBatchProcessor) Stop() error {
        return nil // return nil if not needed
    }
    ```

   > Note: The `Start()` and `Stop()` methods are optional. Implement them if your processing logic requires initialization or cleanup.

2. **Initialize the Batcher**: Create an instance of the `Batcher` by providing the implemented `batchProcessor`, a `batchSize`, and a `flushInterval`:

    ```go
    processor := &MyAwesomeBatchProcessor{}
    batcher := batcher.NewBatcher(processor, 10, 2 * time.Second)
    ```

3. **Add Jobs to the Batcher**:

    ```go
    job := batcher.Job{
        ID:      "unique_job_id",
        Payload: yourPayloadObject, // Replace with your actual payload
    }

    batcher.Add(job)
    ```

4. **Reading the Results**: Ensure to read the results in a non-blocking manner or after all jobs have been added:

    ```go
    for result := range batcher.Results {
        // Handle your job results here
    }
    ```

   > Note: Results will be sent back respectful to the `batchSize` or `flushInterval` that was defined.. Depending on your application's design, consider processing `batcher.Results` in a separate goroutine to avoid blocking if you're adding jobs and reading results concurrently.


### Examples

You can also find a complete usage example in this library.


## License
MIT License
