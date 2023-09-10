<script>
    import { onMount } from "svelte";
    import api from "../api";
    import stores from "../stores";
    import { get } from "svelte/store";
    import Submission from "./Submissions.svelte";

    let signedIn = { username: "", ID: "", authLevel: "" };
    let title = "";

    onMount(async () => {
        title = (await api.GET("/misc/title", {})).body.title;
        let ID = get(stores.user).ID;
        signedIn = (await api.GET("/users/" + ID, {})).body;
    });

    function logout() {
        stores.user.set({ ID: "", token: "", authLevel: "" });
    }
</script>

<h1>{title}</h1>
<h2>{signedIn.username}</h2>
<button on:click={logout}>Logout</button>

<Submission />
