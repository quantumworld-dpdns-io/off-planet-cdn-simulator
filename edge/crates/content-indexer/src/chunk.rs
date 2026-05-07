/// Split text into chunks of at most `max_chars` characters, breaking on word boundaries.
/// Never splits mid-word. If a single word exceeds max_chars, it is included alone as its own chunk.
pub fn chunk_text(text: &str, max_chars: usize) -> Vec<String> {
    let mut chunks: Vec<String> = Vec::new();
    let mut current = String::new();

    for word in text.split_whitespace() {
        if current.is_empty() {
            // Start a new chunk with this word (even if it exceeds max_chars alone)
            current.push_str(word);
        } else {
            // Check if adding a space + word would exceed max_chars
            let needed = current.len() + 1 + word.len();
            if needed > max_chars {
                // Push the current chunk and start fresh
                chunks.push(current.clone());
                current.clear();
                current.push_str(word);
            } else {
                current.push(' ');
                current.push_str(word);
            }
        }
    }

    if !current.is_empty() {
        chunks.push(current);
    }

    // If the input was empty or all whitespace, return a single empty-string chunk
    // to match caller expectations (at least one chunk).
    if chunks.is_empty() {
        chunks.push(String::new());
    }

    chunks
}

#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_basic_chunking() {
        let text = "hello world foo bar baz";
        let chunks = chunk_text(text, 10);
        for chunk in &chunks {
            // No chunk should exceed max_chars unless it's a single oversized word
            assert!(chunk.len() <= 10 || !chunk.contains(' '));
        }
    }

    #[test]
    fn test_single_long_word() {
        let text = "superlongword";
        let chunks = chunk_text(text, 5);
        assert_eq!(chunks, vec!["superlongword"]);
    }

    #[test]
    fn test_empty_text() {
        let chunks = chunk_text("", 512);
        assert_eq!(chunks.len(), 1);
        assert_eq!(chunks[0], "");
    }
}
