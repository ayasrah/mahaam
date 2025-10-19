# Mahaam App

### Overview

Mahaam is an open-sourced, production-ready ToDo app implemented in **`C#, Java, Go, TypeScript, and Python`**.

This documentation explores the following **backend service concepts** through Mahaam.

- **`Service Setup`**: App creation.
- **`Service Design`**: Functions, models, and design.
- **`Service Module`**: Controllers, Services and Repositories.
- **`Service Infra`**: App utilities.
- **`Service Testing`**: Integration tests.

### Purpose

I started Mahaam as a proof of concept for different technologies and architectures. I wanted it to go deeper than the typical shallow ToDo apps, that's why I added real-world functionalites and pushed it to production to complete the cycle.

### Source Code

Currently, Mahaam is implemented in five languages: `C#, Java, Go, TypeScript, and Python` with exact functionality, all expose exact API endpoints and all connected to same database schema.

### Philosophy

Maybe one language is better than another for a specific use case; however, these principles are what Mahaam cares about regardless of language or framework, and they are common throughout:

- Does the app meet **business and user** needs.
- Is the **app model** well designed.
- Is the codebase readable.
- Is the codebase **maintainable**.
- Avoid overengineering (more in Java/C# codebases).
- Avoid spaghetti code (more in JavaScript codebases).

### Target Audience

Mahaam targets software engineers at all levels.

### Explore

Mahaam source code is available on Github and the app is live on the App Store and Play Store.

<div style="display: flex; gap: 20px; align-items: center; flex-wrap: wrap;margin-top: 30px;">
  <a href="https://github.com/ayasrah/mahaam" target="_blank" style="display: inline-flex; align-items: center; text-decoration: none; color: white; background-color:rgb(17, 18, 20); padding: 12px 20px; border-radius: 8px; font-weight: 500; height: 60px; box-sizing: border-box;border: 1px solid #979797;">
    <svg role="img" width="24" height="24" viewBox="0 0 24 24" xmlns="http://www.w3.org/2000/svg" style="margin-right: 8px;" fill="white"><path d="M12 0C5.37 0 0 5.37 0 12c0 5.3 3.438 9.8 8.207 11.385.6.11.793-.26.793-.577v-2.165c-3.338.726-4.042-1.61-4.042-1.61-.546-1.387-1.333-1.756-1.333-1.756-1.09-.745.083-.73.083-.73 1.205.085 1.84 1.24 1.84 1.24 1.07 1.835 2.805 1.305 3.49.998.108-.775.42-1.305.763-1.605-2.665-.3-5.466-1.332-5.466-5.93 0-1.31.467-2.38 1.235-3.22-.123-.303-.535-1.523.117-3.176 0 0 1.008-.322 3.3 1.23a11.5 11.5 0 0 1 3-.405c1.02.005 2.04.137 3 .405 2.29-1.552 3.297-1.23 3.297-1.23.653 1.653.24 2.873.117 3.176.77.84 1.235 1.91 1.235 3.22 0 4.61-2.805 5.625-5.475 5.92.43.37.823 1.1.823 2.22v3.293c0 .32.192.693.8.575C20.565 21.795 24 17.295 24 12c0-6.63-5.37-12-12-12z"/></svg>
    <span>View on GitHub</span>
  </a>
  <a href="https://play.google.com/store/apps/details?id=ayasrah.mahaam" target="_blank">
    <img src="https://upload.wikimedia.org/wikipedia/commons/7/78/Google_Play_Store_badge_EN.svg" alt="Get it on Google Play" style="height: 60px;">
  </a>
  <a href="https://apps.apple.com/us/app/mahaam/id6502533759" target="_blank">
    <img src="https://upload.wikimedia.org/wikipedia/commons/3/3c/Download_on_the_App_Store_Badge.svg" alt="Download on the App Store" style="height: 60px;">
  </a>
</div>

### Sample screens

<div style="display: flex; gap: 20px; align-items: center; flex-wrap: wrap; margin-top:30px;">
  <img src="/plans_screen.jpg" alt="Groups Screen" width="300" style="border: 1px solid #f0f0f0; border-radius:5px;" />
  <img src="/tasks_screen.jpg" alt="Tasks Screen" width="300" style="border: 1px solid #f0f0f0; border-radius:5px;" />
</div>
