# rima-api
API version of the rima app

Getting started on this API for a (much needed) decent rhyming dictionary in Spanish.

It will only have the "consumer" part of the app for now, which simply returns a set of rhymes and synonyms for a specific word, provided we have already analized that word and stored the results in our database.

The analyzing part will come later, when we have finally settled on a final architecture for this.

A lot of TODOs:
* Reading config
* Creating a dependencies object to be passed to controllers
* Add Throttling
* Add Instrumentation
* Strict / exact mode where rhymes are matched syllable by syllable instead of by sound.
