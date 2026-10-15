## Features

- **Issues and Pull Requests**
    - Summarizes open and closed issues and PRs to give an overview of engagement and project progress.
    - Highlights stale issues to prioritize backlog items that need review.

- **Top Reviewers and Mergers**
    - Lists the top reviewers and mergers from LFX Insights, with each person's activity count and share.

- **Project-Specific Data**
    - Tracks stats for both ArgoCD and Argo Rollouts, with options to include other repos.
    - Provides detailed insights to support project-specific contributions.

## Output

The tool prints Markdown you can paste into a Google Doc. With **Tools > Preferences > Enable Markdown** turned on, Docs converts the bullets, bold text, and links. Counts are the links. The output includes:

- Created Issues: Open / Closed, for Argo CD and Argo Rollouts
- Created PRs: Open / Closed, for Argo CD and Argo Rollouts
- Stale Issues, for Argo CD and Argo Rollouts
- Weekly LFX Insights stats: top reviewers and mergers by project, with count and share

![Example](assets/example.png)

## Usage

1. **Build the application:**

   ```bash
   go build -o argo-cd-contrib-insights
    ```
2. **Run the application:**

   ```bash
   ./argo-cd-contrib-insights 
    ```