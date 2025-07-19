# Cimon

## Purpose
The purpose of cimon is to have an interactive ci tool for developers to show their pipelines in a fast manner way without changing tabs between their IDE or developing environment and Gitlab

## Folder structure
In the *cmd* folder we have to **root.go** which is the entry point for cobra. There we have to functions one for executing and on for adding new flags to the rootCmd. In the
utils folder we have different services like the **gitlabService** which provides the functions we call in the **cmd folder files**.