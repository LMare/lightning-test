# Architecture Lightning Playground

## Vue d'ensemble du système

```
┌─────────────────────────────────────────────────────────┐
│                    UTILISATEUR                          │
└────────────────────────┬────────────────────────────────┘
                         │
                         ▼
┌─────────────────────────────────────────────────────────┐
│          FRONTEND (Port 8000)                           │
│  ├─ HTML Templates                                      │
│  ├─ CSS Styling                                         │
│  └─ JavaScript (generic.js)                             │
└────────────────┬────────────────────────────────────────┘
                 │ HTTP/REST
                 ▼
┌─────────────────────────────────────────────────────────┐
│          BACKEND (gRPC)                                 │
│  ├─ lightningHandler (management)                       │
│  ├─ routerHandler (payment routing)                     │
│  ├─ userHandler (user management)                       │
│  └─ streamEventHandler (real-time events)               │
└────────────────┬────────────────────────────────────────┘
                 │
        ┌────────┴─────────┐
        │                  │              
        ▼                  ▼              
┌──────────────┐  ┌──────────────┐  ┌──────────────┐
│  LND-0       │  │  LND-1       │  │  LNDTower    │
│  (Node 1)    │  │  (Node 2)    │  │  (Watchtower)│
│  Port: 9735  │  │  Port: 9735  │  │  Port: 9911  │
│  RPC: 10009  │  │  RPC: 10009  │  │  RPC: 10009  │
└──────┬───────┘  └──────┬───────┘  └───────┬──────┘
       │                 │                  │
       └─────────────────┼──────────────────┘
                         │
                         ▼
                  ┌──────────────┐
                  │     BTCD     │
                  │  Blockchain  │
                  │  (Bitcoin)   │
                  │ RPC: 18556   │
                  │ P2P: 18555   │
                  └──────────────┘
```

## Composants principaux

### 1. **Frontend**
- Serveur web avec interface utilisateur
- Port exposé: `FRONTEND_PORT` (default 8000)
- Communique avec le Backend via gRPC

### 2. **Backend**
- Serveur d'application Go
- Service handlers:
  - `lightningHandler`: Gestion des nœuds Lightning
  - `routerHandler`: Routage des paiements
  - `userHandler`: Gestion utilisateurs
  - `streamEventHandler`: Events en temps réel
- Expose gRPC pour communication interne

### 3. **Bitcoin Network (BTCD)**
- Implémentation blockchain Bitcoin (simnet)
- Simnet: Réseau de simulation Bitcoin pour testing
- RPC: Port 18556 (communication avec LND)
- P2P: Port 18555 (réseau pair-à-pair)
- Mining automatique avec adresse définie

### 4. **Lightning Network Nodes**

#### **LND-0** et **LND-1**
- Nœuds Lightning Network Daemon
- Gestion des canaux de paiement
- Fonctionnalités:
  - **P2P Network**: Port 9735 (connectivité réseau Lightning)
  - **RPC**: Port 10009 (API gRPC)
  - **Watchtower**: Surveillance des canaux (security)
  - **TLS**: Certificats auto-signés

#### **LNDTower** (Watchtower)
- Surveillance des pénalités (penalty transactions)
- Port 9911: Communication avec autres nœuds
- Sécurité: Prévient les fraudes sur les canaux

## Flux de données

### Création d'un canal de paiement
```
LND-0 ──(P2P 9735)──> LND-1
  │                      │
  └──────> BTCD <────────┘
           (RPC 18556)
```

### Paiement sur Lightning
```
User → Frontend → Backend → LND-0 ──(P2P)──> LND-1 → Destination
                    ↓                         ↓
                  Router (payment path)    Router validation
```

## Dépendances
- `LND-0` et `LND-1` dépendent de `BTCD`
- `BTCD` doit être opérationnel avant les nœuds Lightning

## Volumes persistants
| Service | Volume | Chemin |
|---------|--------|--------|
| BTCD | btcd_data | /root/.btcd |
| LND-0 | lnd-0_data | /root/.lnd |
| LND-1 | lnd-1_data | /root/.lnd |
| LNDTower | tower_data | /root/.lnd |

## Configuration réseau (Simnet)
- **Simnet**: Réseau Bitcoin de simulation rapide
- **Génération auto de blocs**: Mining automatique sur BTCD
- **Pas de vraies transactions**: Données de test uniquement
