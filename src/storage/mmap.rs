use memmap2::Mmap;
use std::fs::File;
use std::io;
use std::path::Path;

pub struct MmapFile {
    mmap: Mmap,
}

impl MmapFile {
    pub fn open(path: &Path) -> io::Result<Self> {
        let file = File::open(path)?;
        let mmap = unsafe { Mmap::map(&file)? };

        Ok(Self { mmap })
    }

    pub fn as_slice(&self) -> &[u8] {
        &self.mmap
    }
}
