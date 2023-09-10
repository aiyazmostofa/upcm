<script>
    import { onMount } from "svelte";
    import api from "../api";
    import SubmissionRow from "./SubmissionRow.svelte";
    import time from "../time.js";
    import { get } from "svelte/store";
    import stores from "../stores";

    let problems = [];
    let submissions = [];
    onMount(async () => {
        submissions = (await api.GET("/submissions", {})).body;
        submissions = await fill(submissions);
        submissions = submissions.reverse();
        problems = (await api.GET("/problems", {})).body;
    });

    let error = "";
    async function handleSubmit(event) {
        const form = event.currentTarget;
        let status = 200;
        let response = await fetch(
            new URL(api.BASE + "/submissions?token=" + get(stores.user).token),
            {
                method: "POST",
                body: new FormData(form),
            }
        )
            .then((response) => {
                status = response.status;
                if (response.ok) {
                    return response.json();
                } else {
                    return response.text();
                }
            })
            .then((body) => {
                return {
                    body: body,
                    status: status,
                };
            });
        if (response.status == 200) {
            let submission = response.body;
            submission.problemTitle = (
                await api.GET("/problems/" + submission.problemID, {})
            ).body.title;
            submission.timestamp = new Date(Date.parse(submission.timestamp));
            submission.timestamp = time.format(submission.timestamp);
            submissions = [submission, ...submissions];
            error = "";
        } else {
            error = response.body;
        }
    }

    async function fill(submissions) {
        for (let i = 0; i < submissions.length; i++) {
            let submission = submissions[i];
            submission.username = (
                await api.GET("/users/" + submission.userID, {})
            ).body.username;
            submission.problemTitle = (
                await api.GET("/problems/" + submission.problemID, {})
            ).body.title;
            submission.timestamp = new Date(Date.parse(submission.timestamp));
            submission.timestamp = time.format(submission.timestamp);
        }
        return submissions;
    }

    async function updateSubmissions() {
        submissions = (await api.GET("/submissions", {})).body;
        submissions = await fill(submissions);
        submissions = submissions.reverse();
    }
</script>

<h3>Submit Solution</h3>
<form style="margin-top: 1em;" on:submit|preventDefault={handleSubmit}>
    <label for="problem">Problem: </label>
    <select name="problemID" id="problem">
        {#each problems as problem}
            <option value={problem.ID + ""}>{problem.title}</option>
        {/each}
    </select>
    <br />
    <label for="content">Code File: </label>
    <input type="file" name="content" />
    <br />
    <input type="submit" value="Submit" />
</form>
<p>{error}</p>

<h3>Submissions</h3>

<table>
    <tr>
        <th>Problem</th>
        <th>Timestamp</th>
        <th>Verdict</th>
        <th>Solution</th>
    </tr>
    {#each submissions as submission}
        <SubmissionRow {submission} />
    {/each}
</table>
<br />
<button on:click={updateSubmissions}>Refresh</button>
