<script>
    import api from "../api";

    export let submission;
    let show = false;
    let content = "";
    async function loadSubmission() {
        if (!show) {
            content = (await api.GET("/submissions/" + submission.ID, {})).body
                .content;
        }
        show = !show;
    }
</script>

<tr>
    <td>{submission.ID}</td>
    <td>{submission.username}</td>
    <td>{submission.title}</td>
    <td>{submission.timestamp}</td>
    <td>{submission.verdict}</td>
    <td
        ><button on:click={loadSubmission}>{show ? "Hide" : "Show"}</button>
    </td>
</tr>
<tr>
    <td
        colspan="6"
        style={"text-align: left;" +
            (show ? "display: table-cell;" : "display: none;")}
    >
        <pre>{content}</pre>
    </td>
</tr>
