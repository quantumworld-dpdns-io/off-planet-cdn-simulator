pub fn chunk_text(text: &str, max_tokens: usize) -> Vec<String> {
    let _ = max_tokens; // TODO: implement token-aware chunking (Sprint S8)
    vec![text.to_string()]
}
