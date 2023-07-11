# Commit

![IBFT_prepare](images/IBFT_commit.png)

*Moniz, H. The Istanbul BFT Consensus Algorithm. Algorithm 2. 2020*


## Validation

The paper doesn't include any extra validation for the commit rule. But, as it was with the *Prepare* message, the implementation's *Commit* uses a hash of the value. So, the same hash validation as in the [prepare rule](PREPARE.md) is performed.

## Body (UponCommit)

The only difference in the rule's implementation is:
- upon recieving a quorum of *Commit* messages, an aggregated BLS signature is created. This is also used as return value to the object using the QBFT instance (QBFT controller) and serves as proof of termination.