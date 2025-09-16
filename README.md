# Dagger Cocogitto Module

Reusable dagger module for handling SemVer within a dagger pipeline



---

## Pre-Requisites

- [Dagger](https://dagger.io)
- [Just](https://github.com/casey/just)



---

## Get Started

Make this module part of your dagger module by using

````bash
dagger install https://github.com/stackitcloud/dagger-cocogitto.git
````



---

## Supported Operations

### Callable module functions

- base
- bump
- changelog
- check
- commit
- git-base
- install-hooks
- lob
- verify

The full list of functions including their description can be retrieved by ``dagger functions`` or ``just list``(which calls dagger functions)

There a some optional parameters which can be set depending on the your requirements.
To figure out which params are available for configuration, use ``dagger call [FN_NAME] -h``



---

## Justfile / Environment

The dagger module can be invoked by either calling a function via ``dagger call  [FN_NAME]`` or by calling the just recipe, e.g. ``just get-version``
Just is reading the ``.env`` variable definitions in and passes them to the dagger function. Therefore to be able to use this dagger module with just, the variables defined in ``.env.template`` must be defined with the corresponding values in ``.env``.

You can create a git token in your account settings - Applications/Access Tokens.
The repository url must be written in lower case only e.g. "../Functions/..." will be rejected.



---

## Cog.toml

With the cog.toml additional and replacement behaviour can be defined. 
So you can e.g. rewrite a commonly used ``[skip-ci]`` to ``[ci-skip]`` by specifying in your cog.toml

````bash
skip_ci = "[ci-skip]"
````

Either you use a cog.toml by adding it to your working repository, so that the pipeline discovers the file or 
provide a cog.toml in the bumpSettings for just bumping but to not persisting it

Full Cocogitto ``cog.toml`` config reference:
https://docs.cocogitto.io/reference/config.html#authorsetting