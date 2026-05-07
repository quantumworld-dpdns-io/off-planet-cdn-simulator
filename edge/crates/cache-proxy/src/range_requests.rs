use axum::http::{HeaderMap, StatusCode, header};

pub struct RangeResponse {
    pub status: StatusCode,
    pub body: Vec<u8>,
    pub content_range: Option<String>,
    pub content_length: usize,
}

/// Parses the Range header and slices `data` accordingly.
/// Returns a full 200 response when no Range header is present.
/// Returns 416 Range Not Satisfiable when the range is out of bounds.
pub fn apply_range(data: Vec<u8>, range_header: Option<&str>) -> RangeResponse {
    let total = data.len();

    let Some(range_str) = range_header else {
        return RangeResponse {
            status: StatusCode::OK,
            body: data,
            content_range: None,
            content_length: total,
        };
    };

    // Parse "bytes=start-end" — both start and end are optional
    let range_str = range_str.strip_prefix("bytes=").unwrap_or(range_str);
    let parts: Vec<&str> = range_str.splitn(2, '-').collect();

    let (start, end) = match parts.as_slice() {
        [s, e] => {
            let start: usize = s.trim().parse().unwrap_or(0);
            let end: usize = e.trim().parse().unwrap_or(total.saturating_sub(1));
            (start, end)
        }
        _ => {
            return RangeResponse {
                status: StatusCode::RANGE_NOT_SATISFIABLE,
                body: vec![],
                content_range: Some(format!("bytes */{total}")),
                content_length: 0,
            };
        }
    };

    if start >= total || end >= total || start > end {
        return RangeResponse {
            status: StatusCode::RANGE_NOT_SATISFIABLE,
            body: vec![],
            content_range: Some(format!("bytes */{total}")),
            content_length: 0,
        };
    }

    let slice = data[start..=end].to_vec();
    let len = slice.len();
    RangeResponse {
        status: StatusCode::PARTIAL_CONTENT,
        body: slice,
        content_range: Some(format!("bytes {start}-{end}/{total}")),
        content_length: len,
    }
}

/// Injects Content-Range header into an existing HeaderMap when present.
pub fn inject_range_headers(headers: &mut HeaderMap, rr: &RangeResponse) {
    if let Some(ref cr) = rr.content_range {
        if let Ok(v) = cr.parse() {
            headers.insert(header::CONTENT_RANGE, v);
        }
    }
    if let Ok(v) = rr.content_length.to_string().parse() {
        headers.insert(header::CONTENT_LENGTH, v);
    }
}

#[cfg(test)]
mod tests {
    use super::*;
    use crate::revalidate::make_etag;

    #[test]
    fn full_response_when_no_range() {
        let data = b"hello world".to_vec();
        let rr = apply_range(data.clone(), None);
        assert_eq!(rr.status, StatusCode::OK);
        assert_eq!(rr.body, data);
        assert!(rr.content_range.is_none());
    }

    #[test]
    fn partial_content_for_valid_range() {
        let data = b"hello world".to_vec();
        let rr = apply_range(data, Some("bytes=0-4"));
        assert_eq!(rr.status, StatusCode::PARTIAL_CONTENT);
        assert_eq!(rr.body, b"hello");
        assert_eq!(rr.content_range, Some("bytes 0-4/11".to_string()));
    }

    #[test]
    fn range_not_satisfiable_out_of_bounds() {
        let data = b"hi".to_vec();
        let rr = apply_range(data, Some("bytes=10-20"));
        assert_eq!(rr.status, StatusCode::RANGE_NOT_SATISFIABLE);
    }

    #[test]
    fn etag_deterministic() {
        let e1 = make_etag(1234, 9999);
        let e2 = make_etag(1234, 9999);
        assert_eq!(e1, e2);
    }

    #[test]
    fn etag_differs_on_size_change() {
        let e1 = make_etag(100, 9999);
        let e2 = make_etag(200, 9999);
        assert_ne!(e1, e2);
    }
}
