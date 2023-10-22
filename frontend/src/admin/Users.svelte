<script>
    import { onMount } from "svelte";
    import api from "../api";
    import UserRow from "./UserRow.svelte";

    let users = [];
    onMount(async () => {
        users = (await api.GET("/users", {})).body;
    });

    let username = "";
    let password = "";
    let authLevel = "Team";
    let error = "";
    async function newUser() {
        let response = await api.POST("/users", {
            username,
            password,
            authLevel,
        });
        if (response.status != 200) {
            error = response.body;
            return;
        }
        users = [...users, response.body];
        error = "";
    }
</script>

<main>
    <h3>Users</h3>
    <table>
        <tr>
            <th>Username</th>
            <th>Password</th>
            <th>Auth Level</th>
            <th>Edit</th>
        </tr>
        {#each users as user}
        <UserRow user={user} />
        {/each}
    </table>
    <h3>New User</h3>
    <form on:submit|preventDefault={newUser}>
        <label for="username">Username: </label>
        <input
            type="text"
            name="username"
            id="username"
            bind:value={username}
        />
        <br />

        <label for="password">Password: </label>
        <input
            type="password"
            name="password"
            id="password"
            bind:value={password}
        />
        <br />

        <label for="authLevel">Auth Level: </label>
        <select name="authLevel" id="authLevel" bind:value={authLevel}>
            <option value="Team">Team</option>
        </select>
        <br />

        <input type="submit" value="Create" />
        <br />
    </form>

    <p>{error}</p>
</main>
