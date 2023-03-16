The program initializes a router using the http router package and gets the configuration settings. It creates a PostgreSQL client and initializes a Google translator for translating lessons. It registers the appropriate handlers with the router and starts the HTTP server.
The program provides custom error handling through the apperror package, a middleware function that wraps the appHandler function to handle errors and return appropriate HTTP responses. It also provides repository implementations for Course and Lesson entities that interact with a PostgreSQL database.
Each entity has its own handler, model, service, and storage logic implemented in respective files and directories.
In summary, the program is a web application with flexible and modular code that can be extended or replaced as needed. It has custom error handling and interacts with a PostgreSQL database for storage.

# This is the description of SOLID principles in this project


**Single Responsibility Principle (SRP)** - You have separated the concerns of each component in your project, such as the database layer, handler layer, and service layer, into different packages and interfaces, so each component has a single, well-defined responsibility. For example, the course.handler package is responsible for handling HTTP requests and responses related to course entities, while the course.service package is responsible for providing business logic related to course entities.

**Open/Closed Principle (OCP)** - You have implemented the OCP by using interfaces for certain components in your project, such as the Repository interface for database operations. This allows for the implementation of these components to be easily changed or extended without affecting other parts of the project.

**Liskov Substitution Principle (LSP)** - You have used Go's type system to ensure that objects of different types can be used interchangeably in your code. For example, the course.Service and lesson.Service both implement the Repository interface, so they can be used interchangeably in other parts of the project.

**Interface Segregation Principle (ISP)** - You have used small, specific interfaces in your project, such as the Repository interface, which only defines the methods that are needed for a specific component. This ensures that a component only needs to implement the methods that it needs, rather than a larger, more general interface.

**Dependency Inversion Principle (DIP)** - You have inverted the dependency of high-level components on low-level components by using interfaces and dependency injection. For example, the course.handler package depends on the course.service package, but not the other way around. This allows for the implementation of the course.service package to be changed or extended without affecting the course.handler package.

Here is detailed description for each file

**app.go**
The main function initializes a router using the http router package and gets the configuration settings.
The code creates a PostgreSQL client using the Postrgre SQL package and initializes a Google translator for translating lessons.
It creates repositories for courses and lessons and registers the appropriate handlers with the router.
Finally, the start function listens to the configured socket and starts the HTTP server.

**postgresql.go**
Provides a _PostgreSQL client_ for communicating with a Postgre SQL database. 
The code defines a Client interface that specifies the methods for executing SQL statements.
The NewClient function creates a new PostgreSQL client by taking in the context, maximum number of attempts to connect, and a StorageConfig struct that contains the PostgreSQL database connection details.
It then creates a connection pool with _pgxpool.Connect() method_, using the provided database connection string.
It uses _utils.DoWithTries() function_ to retry the connection process in case of an error.
If the connection is established successfully, the function returns a PostgreSQL connection pool.

**logging.go**
The _writeHook struct_ is defined to define a hook that can be used by the logger to write logs to different outputs (i.e., file and console).
The _Fire method_ of the writeHook struct is used to write the logs to the specified outputs.
The _Levels method_ of the writeHook struct is used to specify the log levels for which the hook should be triggered.
The _Logger struct_ is defined to wrap the _logrus.Entry struct_.
The _GetLoggerWithField method_ of the Logger struct is used to create a new logger with a specified field.
The _init function_ initializes the logger by creating a new logrus instance, setting the log levels, and defining the log output format.
It also creates a new directory for storing log files and creates a new file for writing all logs.
It then sets the writeHook to the logger instance to write logs to the file and console.
Finally, it defines the GetLogger function to return a new Logger instance.

**translation.go and translator.go**

This Go package is responsible for translating lesson names using the Google Cloud API and storing the translations in a PostgreSQL database. 

The _init() function_ creates a new PostgreSQL connection pool by calling postrgresql.NewClient() function and sets the pool to the db variable.
The _TranslateLessonName() function_ takes in a lesson ID, language, and a Translator interface that specifies the method for translating a text to a given language.
It queries the database using the _db.QueryRow() method_ to check if a translation already exists for the specified lesson ID and language.
If the translation already exists, it returns without doing anything.
If the query returns an error other than pgx.ErrNoRows, it returns an error.
If the query returns pgx.ErrNoRows, it calls the _translator.Translate() method_ to get the translated lesson name and stores it in the database using the _db.Exec() method_.
If there is an error while storing the translation in the database, it logs the error and returns without doing anything.
If everything is successful, it returns nil.

**config.go**

This Go package is responsible for reading the application configuration from a YAML file and returning a Config struct that contains the configuration settings. 

The _Config struct_ defines the configuration settings for the application, including IsDebug, Listen, and Storage.
The _StorageConfig struct_ defines the database connection settings for the application.
The _GetConfig() function_ reads the application configuration from a YAML file named config.yml.
It creates a new Config instance and populates it with the values from the YAML file.
If there is an error while reading the YAML file, it logs the error and exits the application.
It uses the sync.Once type to ensure that the configuration is read only once during the lifetime of the application.
It returns the Config instance.

**Error handling in project was done in apperror directory (error.go and middleware.go)**

_error.go_ package defines a custom error type called AppError that can be used to represent application errors. 

The _AppError struct_ defines the error information including the original error, a message for the end user, a message for the developer, and a code that identifies the error.
The ErrNotFound variable is an AppError instance representing a "not found" error. It has a predefined code of "US-0000003".
The _Error() method_ returns the end user message of the AppError instance.
The _Unwrap() method_ returns the original error that caused the AppError.
The _Marshal() method_ returns a JSON representation of the AppError instance.
The _NewAppError() function_ creates a new AppError instance with the provided error, message, developer message, and code.
The _systemError() function_ creates a new AppError instance with a default message and error code, based on the provided error.

_middleware.go_ package defines a middleware function Middleware() that wraps an appHandler function to handle errors and return appropriate HTTP responses. 

The _appHandler function_ is a type that takes an http.ResponseWriter and an http.Request and returns an error.
The _Middleware() function_ takes an appHandler function and returns an http.HandlerFunc.
Inside the _Middleware() function_, an AppError pointer is created to hold any application error that may occur during the execution of the appHandler function.
The _appHandler function_ is called with the given http.ResponseWriter and http.Request objects, and any error returned is captured in the err variable.
If err is not nil, the Content-Type header of the response is set to "application/type".
If err is an instance of AppError, a specific HTTP status code and error message are returned depending on the type of AppError instance.
If err is not an instance of AppError, a generic HTTP status code of 418 (I'm a teapot) and an error message are returned.
The returned error message is serialized into JSON format using the Marshal() method of the appropriate AppError instance or a new systemError() instance.
The _response writer's header_ is written with the correct HTTP status code and the serialized error message is written to the response body.
Overall, this middleware function provides a standard way of handling application errors and returning appropriate HTTP responses.

repository implementation for the Course entity that interacts with a PostgreSQL database was done in **course_db.go**

The _repositoryCourse struct_ defines the repository implementation and holds a postrgresql.Client and a logging.Logger instance.
The _NewRepositoryCourse function_ creates a new repositoryCourse instance and returns it as a course.Repository interface.
The _Create method_ inserts a new course into the database and sets the ID of the provided course. If an error occurs, it logs a detailed error message that includes the PostgreSQL error code and returns the error.
The _FindAll method_ retrieves all courses from the database and returns them as a slice of course.Course instances. If an error occurs, it returns the error.
The _FindOne method_ retrieves a single course from the database based on the provided ID and returns it as a course.Course instance. If an error occurs, it returns the error.
The _Update method_ updates an existing course in the database based on the provided course. If an error occurs, it logs a detailed error message that includes the PostgreSQL error code and returns the error.
The _Delete method_ deletes a course from the database based on the provided ID. If an error occurs, it logs a detailed error message that includes the PostgreSQL error code and returns the error.

repository implementation for the Course entity that interacts with a PostgreSQL database was done in **lesson_db.go**

On the other hand the _repositoryLesson struct_ defines the lesson repository implementation and holds a postrgresql.Client and a logging.Logger instance.
The _NewRepositoryLesson function_ creates a new repositoryLesson instance and returns it as a lesson.Repository interface.
The _FindAll method_ retrieves all lessons from the database and returns them as a slice of lesson.Lesson instances. It also retrieves the courses associated with each lesson and adds them to the respective lesson.Lesson instance. If an error occurs, it returns the error.
The _FindOne method_ retrieves a single lesson from the database based on the provided ID and returns it as a lesson.Lesson instance. It also retrieves the courses associated with the lesson and adds them to the respective lesson.Lesson instance. If an error occurs, it returns the error.
The _Update method_ updates an existing lesson in the database based on the provided lesson. If an error occurs, it logs a detailed error message that includes the PostgreSQL error code and returns the error.

**handler.go** package defines an interface called Handler, which has a single method called Register. The Register method takes an instance of the httprouter.Router type and registers the necessary routes and handlers for the particular HTTP endpoint being implemented. This allows for modular and flexible code that can be easily extended or replaced as needed.

All the **handler, model, service and storage logic for each entity (Course, Lesson and student)** is implemented in respective files and directories.