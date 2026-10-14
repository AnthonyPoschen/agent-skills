pub fn bad(x: u64, v: Option<u8>) -> u32 {
    let _h = std::env::var("HOME");
    std::thread::sleep(std::time::Duration::from_millis(5));
    println!("hi");
    let _ = v.unwrap();
    if x == 0 { panic!("zero"); }
    x as u32
}
pub fn good(x: u64) -> Result<u32, std::num::TryFromIntError> { u32::try_from(x) }
#[cfg(test)]
mod tests { #[test] fn t() { Some(1).unwrap(); println!("ok"); } }
