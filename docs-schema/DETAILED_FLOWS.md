# Schémas Détaillés - Lightning Playground

## 1. Flux d'ouverture d'un canal de paiement

```
Timeline: Création de canal LND-0 → LND-1 (1 BTC)

┌─────────────────────────────────────────────────────────────┐
│                   FRONTEND                                  │
│          User clicks: "Open Channel"                         │
│          Amount: 1 BTC                                       │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                   BACKEND                                   │
│      Call: LND-0.OpenChannel(                               │
│        nodeID: lnd1-pubkey,                                 │
│        amount: 1 BTC,                                       │
│        pushAmount: 0)                                       │
└─────────────────┬───────────────────────────────────────────┘
                  │ gRPC
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                   LND-0                                     │
│  1. Get balance from wallet                                 │
│  2. Create funding transaction (1 BTC)                      │
│  3. Connect to LND-1 (P2P :9735)                            │
│  4. Exchange channel parameters                             │
│  5. Sign transaction                                        │
│  6. Broadcast to BTCD                                       │
└─────────────────┬───────────────────────────────────────────┘
                  │ RPC :18556
                  ▼
┌─────────────────────────────────────────────────────────────┐
│                   BTCD                                      │
│  1. Validate transaction                                    │
│  2. Accept into mempool                                     │
│  3. Mine in next block (~simnet instant)                    │
│  4. Broadcast to network                                    │
└─────────────────┬───────────────────────────────────────────┘
                  │ P2P :18555
                  ▼
┌─────────────────────────────────────────────────────────────┐
│               LND-0 & LND-1                                 │
│  1. Both watch BTCD for confirmation                        │
│  2. 1st confirmation: Channel tentative                     │
│  3. 6 confirmations: Channel active                         │
└─────────────────┬───────────────────────────────────────────┘
                  │
                  ▼
┌─────────────────────────────────────────────────────────────┐
│           CHANNEL OPERATIONAL ✓                             │
│                                                              │
│  LND-0 Side    │ LND-1 Side                                 │
│  Balance: 1    │ Balance: 0 (peut recevoir)                │
│  Reserved: 1   │ Reserved: 1 (locked in channel)            │
└─────────────────────────────────────────────────────────────┘
```

---

## 2. Flux de paiement Lightning (Payment)

```
Scénario: LND-0 envoie 0.1 BTC à LND-1

┌────────────────────────────────────────┐
│    LND-0 Wallet Balance: 1.0 BTC       │
│    Channel LND-0↔LND-1: 1.0 vs 0.0     │
│                                        │
│    LND-1 Wallet Balance: 0 BTC         │
│    Channel LND-0↔LND-1: 0.0 vs 1.0     │
└────────────────────────────────────────┘

Step 1: Create Payment Request (LND-1)
┌────────────────────────────────────────┐
│ Backend calls LND-1.AddInvoice:        │
│  - Amount: 0.1 BTC                     │
│  - Memo: "Payment for service"         │
│  - Expiry: 1 hour                      │
│                                        │
│ LND-1 returns: payment_request (LN URI)│
│ Format: lnbc100m1...@lnd1              │
└────────────────────────────────────────┘
                  │
                  ▼
Step 2: Send Payment (LND-0)
┌────────────────────────────────────────┐
│ Backend calls LND-0.SendPayment:       │
│  - PaymentRequest: lnbc100m1...        │
│  - Amount: 0.1 BTC (converted to msat) │
│                                        │
│ LND-0 does route finding:              │
│  → Direct route: LND-0 → LND-1 ✓       │
│    Fees: 0 (direct, no intermediary)   │
└────────────────────────────────────────┘
                  │
                  ▼
Step 3: Execute Payment
┌────────────────────────────────────────┐
│ Payment Flow:                          │
│                                        │
│ LND-0           LND-1                  │
│ 1.0 BTC   →    0.0 BTC                 │
│                                        │
│ Decrease 0.1  Increase 0.1             │
│                                        │
│ 0.9 BTC   →    0.1 BTC                 │
│                                        │
│ (No blockchain transaction!)           │
│ (Updated state commitment only)        │
└────────────────────────────────────────┘
                  │
                  ▼
Step 4: Payment Confirmation
┌────────────────────────────────────────┐
│ Both nodes update local commitment tx  │
│ Hash: preimage revealed                │
│ HTLC (Hash Time Locked Contract) ✓     │
│                                        │
│ Frontend gets: Success                 │
└────────────────────────────────────────┘

Final State:
┌────────────────────────────────────────┐
│    LND-0 Wallet Balance: 0.9 BTC       │
│    Channel LND-0↔LND-1: 0.9 vs 0.1     │
│                                        │
│    LND-1 Wallet Balance: 0.1 BTC       │
│    Channel LND-0↔LND-1: 0.1 vs 0.9     │
└────────────────────────────────────────┘
```

---

## 3. Flux Watchtower - Détection de fraude

```
Scénario: LND-0 tente de fermer le canal avec état périmé

NORMAL SCENARIO (Honest closure):
┌─────────────────────────────────────────────┐
│ LND-0 & LND-1 s'accordent sur état final   │
│ Broadcast cooperative close transaction    │
│ → No watchtower action needed               │
└─────────────────────────────────────────────┘

FRAUD SCENARIO (Invalid state closure):
┌──────────────────────────────────────────────────────┐
│          State Evolution                             │
├──────────────────────────────────────────────────────┤
│ State 1: LND-0 = 1.0 BTC, LND-1 = 0.0 BTC           │
│ State 2: LND-0 = 0.9 BTC, LND-1 = 0.1 BTC           │
│ State 3: LND-0 = 0.8 BTC, LND-1 = 0.2 BTC ← Latest │
│ State 2: LND-0 = 0.9 BTC, LND-1 = 0.1 BTC ← OLD!    │
└──────────────────────────────────────────────────────┘
                  │
                  ▼
LND-0 (offline) broadcasts State 2 commitment
                  │
                  ▼
┌─────────────────────────────────────────────────────┐
│           Transaction in mempool                    │
│         LND-0 trying to get 0.9 BTC                 │
│         (vs actual 0.8 BTC in State 3)              │
└─────────────────────────────────────────────────────┘
                  │
                  ▼
         LNDTower Watchtower
         monitors BTCD:9911
                  │
                  ▼
    ✓ Detected: Invalid state!
    ✓ State 2 < State 3
                  │
                  ▼
    Watchtower creates Penalty TX:
    - Takes ALL 1 BTC from channel
    - Sends to LND-1 as punishment
                  │
                  ▼
    Broadcasts penalty to BTCD
                  │
                  ▼
LND-0's fraud tx rejected,
Penalty tx confirms first
                  │
                  ▼
    LND-1 receives 1 BTC
    LND-0 loses 1 BTC (complete channel)
    
    RESULT: Honest node protected! ✓
```

---

## 4. Data Flow Architecture

```
┌──────────────────────────────────────────────────────────┐
│                    USER LAYER                            │
│              (Frontend Browser)                          │
└─────────────────────────┬────────────────────────────────┘
                          │ HTTP/REST
                          │ (JSON payloads)
                          ▼
┌──────────────────────────────────────────────────────────┐
│                   APPLICATION LAYER                      │
│                    (Go Backend)                          │
│                                                           │
│  ┌──────────────────────────────────────────────────┐  │
│  │ Handler Functions:                               │  │
│  │ • lightningHandler.GetNodeInfo()                │  │
│  │ • lightningHandler.OpenChannel()                │  │
│  │ • routerHandler.SendPayment()                   │  │
│  │ • userHandler.GetUser()                         │  │
│  │ • streamEventHandler.SubscribePayments()        │  │
│  └──────────────────────────────────────────────────┘  │
└─────────┬────────────────────────────────────────────────┘
          │ gRPC with mTLS
          │ protobuf serialization
          │
    ┌─────┴─────────┐
    │               │
    ▼               ▼
 ┌─────┐        ┌─────┐
 │ LND0│        │ LND1│
 └──┬──┘        └──┬──┘
    │              │
    │ P2P:9735     │
    ├──────────────┤
    │              │
    └──────────────┘
            │
            │ Both LND nodes register with Watchtower
            │ (P2P connection, not Backend)
            │
            ▼
        ┌─────┐
        │TOWER│
        └─────┘
    
    (Lightning Network)
            │
            │ RPC:18556 (Bitcoin protocol)
            ▼
         ┌─────┐
         │BTCD │
         └─────┘
```

---

## 5. Channel State Representation

```
Visual representation of channel balance

Initial state (1 BTC locked):
┌────────────────────────────────────┐
│ LND-0 │                  │ LND-1    │
│ 1.000 │ Reserved: 1 BTC  │ 0.000    │
└────────────────────────────────────┘

After 0.1 BTC payment (LND-0 → LND-1):
┌────────────────────────────────────┐
│ LND-0 │                  │ LND-1    │
│ 0.900 │ Commitment tx    │ 0.100    │
└────────────────────────────────────┘
       ↓                          ↓
    Can send more          Can receive more
    up to 0.900            up to 0.900

After 0.05 BTC payment (LND-1 → LND-0):
┌────────────────────────────────────┐
│ LND-0 │                  │ LND-1    │
│ 0.950 │ Updated commit   │ 0.050    │
└────────────────────────────────────┘

Max capacity: Always = 1.000 BTC
```

---

## 6. Error & Recovery Scenarios

```
┌───────────────────────────────────────────────────┐
│             CHANNEL CLOSURE TYPES                 │
├───────────────────────────────────────────────────┤

1. COOPERATIVE CLOSE (Happy path)
   ┌─────────────────────────────┐
   │ Both nodes agree on final    │
   │ state and signing            │
   └────────────┬────────────────┘
                │
                ▼
        Both broadcast cooperatively
        → Minimal fees (~1 sat/byte)
        → Fast settlement

2. FORCE CLOSE (One party offline)
   ┌─────────────────────────────┐
   │ Node needs to close channel  │
   │ counterparty is offline      │
   └────────────┬────────────────┘
                │
                ▼
        Broadcast latest commitment
        → Watchtower watches for fraud
        → Higher fees
        → Longer wait (~144 blocks for CSV)

3. BREACH (Fraud detected)
   ┌─────────────────────────────┐
   │ Node broadcasts old state    │
   │ Watchtower detects           │
   └────────────┬────────────────┘
                │
                ▼
        Penalty transaction triggered
        → Fraudster loses entire channel
        → Honest party gets all funds

```
