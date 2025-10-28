# rr — A small Go CLI to rename episode files

English | [简体中文](./README_zh.md)

rr is a lightweight command-line tool written in Go that scans a directory for video and subtitle files, extracts episode numbers and bracketed tags from filenames, and renames files into a standardized "Show SxxExx [tags].ext" format. It supports a dry-run mode, preserves existing bracketed tags by default, and provides colored terminal output (can be disabled).
