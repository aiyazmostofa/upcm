<script>
    import { onMount } from "svelte";
    import time from "../time";
    import api from "../api";

    let dryRunStartTime = "";
    let contestStartTime = "";
    let contestEndTime = "";

    let status = "";

    onMount(async () => {
        let times = (await api.GET("/misc/times", {})).body;
        dryRunStartTime = time.format(
            new Date(Date.parse(times.dryRunStartTime))
        );
        contestStartTime = time.format(
            new Date(Date.parse(times.contestStartTime))
        );
        contestEndTime = time.format(
            new Date(Date.parse(times.contestEndTime))
        );
    });

    async function setTimes() {
        let response = await api.PUT("/misc/times", {
            dryRunStartTime,
            contestStartTime,
            contestEndTime,
        });
        if (response.status != 200) {
            status = response.body;
            return;
        } else {
            status = "times saved";
        }

        let times = response.body;
        dryRunStartTime = time.format(
            new Date(Date.parse(times.dryRunStartTime))
        );
        contestStartTime = time.format(
            new Date(Date.parse(times.contestStartTime))
        );
        contestEndTime = time.format(
            new Date(Date.parse(times.contestEndTime))
        );
    }
</script>

<h3>Settings</h3>
<p>Format: MM/DD/YYYY|HH:MM:SS|TZ</p>

<form on:submit|preventDefault={setTimes}>
    <label for="dry">Dry Run Start Time: </label>
    <input type="text" name="dry" id="dry" bind:value={dryRunStartTime} />
    <br />

    <label for="wet">Contest Start Time: </label>
    <input type="text" name="wet" id="wet" bind:value={contestStartTime} />
    <br />

    <label for="end">Contest End Time: </label>
    <input type="text" name="end" id="end" bind:value={contestEndTime} />
    <br />

    <input type="submit" value="Save" />
</form>

<p>{status}</p>
