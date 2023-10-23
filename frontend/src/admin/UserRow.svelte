<script>
    import api from "../api";

    export let user;
    let edit = false;

    let backupUsername;
    let backupPassword;

    let message = "";

    let here = true;

    async function editUser() {
        if (edit) {
            let response = await api.PUT("/users/" + user.ID, {
                username: user.username,
                password: user.password,
            });

            if (response.status != 200) {
                message = response.body;
                user.username = backupUsername;
                user.password = backupPassword;
                return;
            }
        } else {
            backupUsername = user.username;
            backupPassword = user.password;
        }
        edit = !edit;
    }

    async function deleteUser() {
        let response = await api.DELETE("/users/" + user.ID, {});
        if (response.status != 200) {
            message = response.body;
            return;
        }
        here = false;
    }
</script>

<tr style={here ? "" : "display: none;"}>
    {#if edit}
        <td
            ><input
                style="width: 6rem"
                type="text"
                name="username"
                bind:value={user.username}
            /></td
        >
        <td
            ><input
                style="width: 6rem"
                type="text"
                name="password"
                bind:value={user.password}
            /></td
        >
        <td>{user.authLevel}</td>
        <td><button on:click={editUser}>Save</button></td>
        <td><button on:click={deleteUser}>Delete</button></td>
        <td>{message}</td>
    {:else}
        <td>{user.username}</td>
        <td>{user.password}</td>
        <td>{user.authLevel}</td>
        <td><button on:click={editUser}>Yes</button></td>
    {/if}
</tr>
