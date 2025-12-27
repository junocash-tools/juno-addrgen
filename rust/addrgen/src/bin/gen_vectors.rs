use orchard::keys::{FullViewingKey, Scope, SpendingKey};
use serde::Serialize;
use zip32::AccountId;

use juno_addrgen::{HRP_JUNO_UA, HRP_JUNO_UFVK, JUNO_COIN_TYPE, TYPECODE_ORCHARD};

#[derive(Serialize)]
struct VectorsV1 {
    version: u32,
    ufvk: String,
    addresses: Vec<String>,
}

fn main() {
    let seed = [7u8; 64];
    let account = AccountId::try_from(0).expect("account");
    let sk = SpendingKey::from_zip32_seed(&seed, JUNO_COIN_TYPE, account).expect("sk");
    let fvk = FullViewingKey::from(&sk);

    let ufvk = juno_addrgen::zip316::encode_unified_container(
        HRP_JUNO_UFVK,
        TYPECODE_ORCHARD,
        &fvk.to_bytes(),
    )
    .expect("ufvk");

    let mut addresses = Vec::with_capacity(100);
    for index in 0u32..100u32 {
        let raw = fvk.address_at(index, Scope::External).to_raw_address_bytes();
        let addr = juno_addrgen::zip316::encode_unified_container(
            HRP_JUNO_UA,
            TYPECODE_ORCHARD,
            &raw,
        )
        .expect("addr");
        addresses.push(addr);
    }

    let v = VectorsV1 {
        version: 1,
        ufvk,
        addresses,
    };
    println!("{}", serde_json::to_string_pretty(&v).expect("json"));
}

