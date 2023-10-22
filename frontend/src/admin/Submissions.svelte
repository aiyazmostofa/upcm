<script>
  import { onMount } from "svelte";
  import api from "../api";
  import SubmissionRow from "./SubmissionRow.svelte";
  import time from "../time.js";

  let submissions = [];

  async function updateSubmissions() {
    submissions = (await api.GET("/submissions/best", {})).body;
    scores = await calculateScores(submissions);
    submissions = await fill(submissions);
    submissions = submissions.reverse();
  }

  onMount(async () => {
    updateSubmissions();
  });

  async function fill(submissions) {
    for (let i = 0; i < submissions.length; i++) {
      let submission = submissions[i];
      submission.timestamp = new Date(Date.parse(submission.timestamp));
      submission.timestamp = time.format(submission.timestamp);
    }
    return submissions;
  }

  let scores = [];

  async function calculateScores(submissions) {
    let attempts = {};
    let scores = {};
    for (let i = 0; i < submissions.length; i++) {
      let submission = submissions[i];
      if (submission.verdict == "Internal Error" || submission.problemID == 1)
        continue;
      if (!attempts.hasOwnProperty(submission.username)) {
        attempts[submission.username] = {};
      }

      if (!attempts[submission.username].hasOwnProperty(submission.problemID)) {
        attempts[submission.username][submission.problemID] = {
          right: false,
          wrong: 0,
        };
      }

      let attempt = attempts[submission.username][submission.problemID];
      if (attempt.right) continue;
      if (submission.verdict == "Correct Answer") {
        if (!scores.hasOwnProperty(submission.username)) {
          scores[submission.username] = 60;
        } else {
          scores[submission.username] += 60;
        }
        scores[submission.username] -= Math.min(5 * attempt.wrong, 60);
        attempt.right = true;
      } else {
        attempt.wrong++;
      }
    }

    let users = (await api.GET("/users", {})).body;
    for (let i = 0; i < users.length; i++) {
      if (
        !scores.hasOwnProperty(users[i].username) &&
        users[i].authLevel == "Team"
      ) {
        scores[users[i].username] = 0;
      }
    }

    let result = Object.entries(scores);
    result.sort((a, b) => {
      if (a[1] == b[1]) {
        return 0;
      }
      return a[1] < b[1] ? 1 : -1;
    });
    return result;
  }
</script>

<h3>Scores</h3>
<table>
  <tr><th>Team</th><th>Score</th></tr>
  {#each scores as score}
    <tr><td>{score[0]}</td><td>{score[1]}</td></tr>
  {/each}
</table>
<br />
<button on:click={updateSubmissions}>Refresh Both Tables</button>

<h3>Submissions</h3>

<table>
  <tr>
    <th>User</th>
    <th>Problem</th>
    <th>Timestamp</th>
    <th>Verdict</th>
    <th>Solution</th>
  </tr>
  {#each submissions as submission}
    <SubmissionRow bind:submission />
  {/each}
</table>
<br />
