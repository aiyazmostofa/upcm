<script>
    import stores from "./stores";
    import api from "./api";

    let username = "";
    let password = "";
    let error = "";

    async function login() {
        let response = await api.GET("/token", { username, password });

        if (response.status == 200) {
            stores.user.set(response.body);
        } else {
            error = response.body;
        }
    }
</script>

<h1>UIL Programming Contest Manager</h1>
<h3>Login</h3>

<form on:submit|preventDefault={login}>
    <label for="username">Username: </label>
    <input type="text" id="username" name="username" bind:value={username} />
    <br />

    <label for="password">Password: </label>
    <input
        type="password"
        id="password"
        name="password"
        bind:value={password}
    />
    <br />

    <input type="submit" value="Login" />
</form>
<p>{error}</p>
