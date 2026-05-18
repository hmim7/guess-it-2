#!/bin/sh
# This script runs the student's Go program for the guess-it-2 project.
# We set the "median" strategy as the default for maximum resilience.
export PREDICT_CENTER="median"
./student/guess-it-2
