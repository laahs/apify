# apify
The aim of this package is to define boilerplate to ease api definition. Given a model, a set of methods allowed and a set of restrictions the framework should generate routes/handlers, validate data and handle requests. It should be able to handle models related to other models.

It should be able to implement:
Rate limiters enforce upper thresholds on incoming request throughput.
Serialization is the conversion of language specific data structures to a byte stream for presentation to another system. That other system is commonly a browser (json/xml/html) or a database, among others.
Logging is the time-ordered, preferably structured, output from an application and its constituent components.
Metrics are a record of the instrumented parts of your application and includes the aggregated measurements of latency, request counts, health and others.
Circuit breakers prevent thundering herds thereby improving resiliency against intermittent errors.
Request tracing across multiple services is an important tool for diagnosing issues and recreating the state of the system as a whole.
Probably out of scope: Service discovery allows different services to find each other given known, stable names and the realities of the cloud, where individual systems come and go dynamically and when least expected.


# API flow definition
General flow through an api is the following, each step is pass or fail and return:

- ## Step 0: Router directs the data to the appropriate Handler
 --> return ERROR "no content" if path not allowed, ERROR "not allowed" if method not allowed on content/path.

ACTORS:
- ## Step 1: Initiator goes through a set of middleware to initiate basic info about the query (store query timestamp, parameters, set timeouts, language/country etc in the context for further analysis, metrics, instrumentation, log, debug, rate limiting -> use context)
 --> store in context to passe along with the request to next actor, write method for some objects as log.

- ## Step 2: Checker evaluates header and query parameters, submits data to the "barrier" (a set of middlewares checking authorization, authentity etc)
 --> return ERROR "not authorized" if does not pass authorizations middleware (Basic auth, api access auth (jwt), content access auth (jwt)), ERROR "bad request" if expected headers have not been sent or wrong content-type, ERROR "timeout gateway" if the process took too long
 
- ## Step 3: Extractor sucks the data from request body and returns a payload map[string]interface{}
 --> return ERROR "bad request" if wrong content-type detected/read, ERROR "Internal server error" if the data could not be parsed

- ## Step 4: Evaluator submits payload for evaluation: reject unexpected fields (by name) based on a set of permitted fields from schema/item, reject unvalid data based on a set of constraints for each "field", normalize data (set defaults, format, set private fields values depending on action requested (create, add, remove, replace)). Fields should be accessible through field/subfield or field.subfield pattern.
 --> return ERROR "bad request" if some fields are not valid (unexpected names), ERROR "bad request" if some fields can not be defined/updated, ERROR "bad request" if wrong value types, ERROR "bad request" if values do not fulfill requirements 

- ## Step 5: Generator generates the item/object to be stored from the normalized validated payload
 --> return ERROR "internal server error" if problem
 
- ## Step 6: Dealer manipulates database depending on action expected with item -> Insert, Delete, Update, Patch, Find, FindList (=search)
 --> return ERROR "internal server error" if action could ne be performed

WRAPPER:
- ## Step 7: Responder generates response (json format) with standard format: response object defined as { data map[string]interface{}, statuscode int32, errors map[string]string matching issuer:issue} from the set of inputs (format string, data interface{}, language string, statuscode int32, errors map[string]string). Set response header depending on status code. Write body with response object depending on format. Responder should wrap all other Actors of the Flow.
--> return http response
