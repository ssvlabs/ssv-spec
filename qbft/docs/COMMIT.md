# Commit

![IBFT_prepare](images/IBFT_commit.png)

*Moniz, H. The Istanbul BFT Consensus Algorithm. Algorithm 2. 2020*


## Validation

The paper doesn't include any extra validation for the commit rule. But, as it was with the *Prepare* message, the implementation's *Commit* uses a hash of the value. So, the same hash validation as in the [prepare rule](PREPARE.md) is performed.

## Body (UponCommit)

The differences in the rule's implementation are:
- upon recieving a quorum of *Commit* messages, an aggregated BLS signature is created. This is used as return value to the object using the QBFT instance (QBFT controller) and serves as proof of termination.
- The timer is not set to stop. However, after the instance is decided no *Round-Changes* are processed.