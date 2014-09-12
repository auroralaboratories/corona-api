# Corona API

The Corona API is a REST API server designed to provide access to common desktop interactions, data, and functionality used by the programs that comprise a Linux desktop environment.  With Corona API, desktop applications like taskbars, dashboards, widgets, and launchers can be built using modern web technologies (HTML, Javascript, CSS, etc.).  This offers the possibility of creating highly-interactive user experiences using extremely common tools and languages.

# Why Would Anyone Do This?

Today's graphical desktops are getting more and more complex; the de facto idioms of how we interact with computers being challenged every day.  As platforms like Windows(TM) and Apple(R) OS X(TM) continue to evolve and refine the next generation of standard interface design, the open source community is also moving to keep up.  This has led to the creation of many new tools and libraries for pushing the modern Linux desktop forward.  However, many of these tools are built on the foundation of decades-old technologies that, in many cases, struggle to keep up.  The result is that it takes a lot more work to build unconventional, non-idiomatic interfaces and applications than it does to build ones that the toolkits were designed for and are expecting.

Web applications, however, have always had a need to be extremely flexible and unconventional in their interface design.  The tools and languages that make up those applications have, in turn, also enabled a much broader and more expressive toolkit for building highly custom user experiences.  Corona aims to extend that flexibility onto the modern Linux desktop.


# If My Launcher Returns a `404` Error I Will Lose My Mind

That's a totally fair concern.  Like anything new entering an old arena, there will need to be a shift in how common development practices are performed.  Compiled applications have the very positive trait of being somewhat more predictable in how they behave, with extremely well-known failure modes that people understand how to work with (and work around).  Additionally, in building a web application designed to be run outside of the typical browser environment, the expectations of that application change dramatically.  

That said, modern web technologies provide a plethora of options to make interacting with them consistent with the expectations of a traditional desktop application.  This project will provide reference implementations of several common desktop applications to serve as an example of how these goals can be achieved.  These examples can be found [here](https://github.com/auroralaboratories/corona-ui).
