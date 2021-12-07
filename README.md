Dolt
===

This is a simple yet functional payload service, that meant to mock real world 
service behavior. 

Abstract
---

Any service can be unpredictably slow in a startup or shutdown processes, as 
well as, it can accidentally exit or, in opposite, resist to become terminated 
from external command.

Moreover, such real-world services, can easily exceed any limits, like CPU time 
or granted RAM.

Usage
---

To control application's behaviour, you need to pass one or several 
_Environment_ variables from this list:

* EXITCODE (int), — to control with which exit code the application will use;
* HEALTHURI (string, urlformated), — to control on which _host_, _port_ health 
web server will started, and by which _path_ it will be answering;
* IGNORESIGS (string, comma separated), — defines which os Signals should be 
ignored **KILL (9)** will never be ignored or catched;
* INITTIME (duration: N[ns|us|ms|s|m|h] or their combination), — defines for 
how long it will takes to initialize the application;
* LIFETIME (duration: N[ns|us|ms|s|m|h] or their combination), — defines life 
time of the application, after which it will be stopped;
* STOPTIME (duration: N[ns|us|ms|s|m|h] or their combination), — defines for 
how long it will takes to stop the application;
