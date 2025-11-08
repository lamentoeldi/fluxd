# Flux

Fluxd is a lightweight task scheduler daemon inspired by Airflow.  
It manages and executes tasks, supports dependencies between tasks, ensures fault-tolerance, and logs task execution.

---

## Table of Contents

- [Overview](#overview)
- [Key Features](#key-features)
- [License](LICENSE)

---

## Overview

fluxd is a daemon that runs tasks according to a schedule and respects dependencies between them.  
It is designed to automate workflows, handle failures gracefully, and keep a reliable history of task executions.

---

## Key Features

- Execute tasks based on scheduled times (cron-like expressions).
- Support dependencies between tasks.
- Retry failed tasks and enforce execution timeouts.
- Record task execution logs and history.
- Maintain task metadata in a persistent store.

---
