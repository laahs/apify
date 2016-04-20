# apify
The aim of this package is to define boilerplate to ease api definition. Given a model, a set of methods allowed and a set of restrictions the framework should generate routes/handlers, validate data and handle requests. It should be able to handle models related to other models.

# API flow definition
General flow through an api is the following, each step is pass or fail and return:

0) Router directs the data to the appropriate Handler
 --> return ERROR "no content" if path not allowed, ERROR "not allowed" if method not allowed on content/path

ACTORS:
1) Initiator goes through a set of middleware to initiate basic info about the query (store query timestamp, parameters, set timeouts, language/country etc in the context for further analysis, metrics, instrumentation, log, debug, rate limiting -> use context)

2) Checker evaluates header and query parameters, submits data to the "barrier" (a set of middlewares checking authorization, authentity etc)
 --> return ERROR "not authorized" if does not pass authorizations middleware (Basic auth, api access auth (jwt), content access auth (jwt)), ERROR "bad request" if expected headers have not been sent or wrong content-type, ERROR "timeout gateway" if the process took too long
 
3) Extractor sucks the data from request body and returns a payload map[string]interface{}
 --> return ERROR "bad request" if wrong content-type detected/read, ERROR "Internal server error" if the data could not be parsed

4) Evaluator submits payload for evaluation: reject unexpected fields (by name) based on a set of permitted fields from schema/item, reject unvalid data based on a set of constraints for each "field", normalize data (set defaults, format, set private fields values depending on action requested (create, add, remove, replace)). Fields should be accessible through field/subfield or field.subfield pattern.
 --> return ERROR "bad request" if some fields are not valid (unexpected names), ERROR "bad request" if some fields can not be defined/updated, ERROR "bad request" if wrong value types, ERROR "bad request" if values do not fulfill requirements 

5) Generator generates the item/object to be stored from the normalized validated payload
 --> return ERROR "internal server error" if problem
 
6) Dealer manipulates database depending on action expected with item -> Insert, Delete, Update, Patch, Find, FindList (=search)
 --> return ERROR "internal server error" if action could ne be performed

WRAPPER:
7) Responder generates response (json format) with standard format: response object defined as { data map[string]interface{}, statuscode int32, errors map[string]string matching issuer:issue} from the set of inputs (format string, data interface{}, language string, statuscode int32, errors map[string]string). Set response header depending on status code. Write body with response object depending on format. Responder should wrap all other Actors of the Flow.
--> return http response
