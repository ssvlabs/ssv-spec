package types

import (
	"encoding/json"
	"fmt"
	"github.com/bloxapp/ssv-spec/qbft"
	"github.com/golang/snappy"
	"github.com/klauspost/compress/s2"
	"github.com/stretchr/testify/require"
	"testing"

	v1 "github.com/attestantio/go-eth2-client/api/v1"
	"github.com/attestantio/go-eth2-client/spec/altair"
	"github.com/attestantio/go-eth2-client/spec/phase0"
)

var result []byte

var cd2 = &ConsensusData{
	Duty: &Duty{
		Type: 4,
		PubKey: phase0.BLSPubKey{
			0,
			1,
			2,
			3,
			4,
			5,
			6,
			7,
			8,
			9,
			10,
			11,
			12,
			13,
			14,
			15,
			16,
			17,
			18,
			19,
			20,
			21,
			22,
			23,
			24,
			25,
			26,
			27,
			28,
			29,
			30,
			31,
			32,
			33,
			34,
			35,
			36,
			37,
			38,
			39,
			40,
			41,
			42,
			43,
			44,
			45,
			46,
			47,
		},
		Slot:                    1,
		ValidatorIndex:          2,
		CommitteeIndex:          3,
		CommitteeLength:         128,
		CommitteesAtSlot:        4,
		ValidatorCommitteeIndex: 61,
	},
	AttestationData: &phase0.AttestationData{
		Slot:  100,
		Index: 1,
		BeaconBlockRoot: phase0.Root{
			0,
			1,
			2,
			3,
			4,
			5,
			6,
			7,
			8,
			9,
			10,
			11,
			12,
			13,
			14,
			15,
			16,
			17,
			18,
			19,
			20,
			21,
			22,
			23,
			24,
			25,
			26,
			27,
			28,
			29,
			30,
			31,
		},
		Source: &phase0.Checkpoint{
			Epoch: 1,
			Root: phase0.Root{
				32,
				33,
				34,
				35,
				36,
				37,
				38,
				39,
				40,
				41,
				42,
				43,
				44,
				45,
				46,
				47,
				48,
				49,
				50,
				51,
				52,
				53,
				54,
				55,
				56,
				57,
				58,
				59,
				60,
				61,
				62,
				63,
			},
		},
		Target: &phase0.Checkpoint{
			Epoch: 2,
			Root: phase0.Root{
				64,
				65,
				66,
				67,
				68,
				69,
				70,
				71,
				72,
				73,
				74,
				75,
				76,
				77,
				78,
				79,
				80,
				81,
				82,
				83,
				84,
				85,
				86,
				87,
				88,
				89,
				90,
				91,
				92,
				93,
				94,
				95,
			},
		},
	},
}

func testSSVMsgSSZ() []byte {
	cdSSZ, _ := cd2.MarshalSSZ()
	proposalData := &qbft.ProposalData{
		Data:                     cdSSZ,
		RoundChangeJustification: []*qbft.SignedMessage{},
		PrepareJustification:     []*qbft.SignedMessage{},
	}

	pdSSZ, _ := proposalData.MarshalSSZ()

	msg := &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     1,
		Round:      2,
		Identifier: []byte{1, 2, 3, 4},
		Data:       pdSSZ,
	}

	signedMsg := &qbft.SignedMessage{
		Signature: []byte{1, 2, 3, 4},
		Signers:   []OperatorID{1, 2, 3, 4},
		Message:   msg,
	}

	signedMsgSSZ, _ := signedMsg.MarshalSSZ()
	signedMsgSSZS2 := s2.Encode([]byte{}, signedMsgSSZ)

	msgID := MessageID{}
	copy(msgID[:], msg.Identifier)

	msgToBroadcast := &SSVMessage{
		MsgType: SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    signedMsgSSZS2,
	}
	msgToBroadcastSSZ, _ := msgToBroadcast.MarshalSSZ()
	return msgToBroadcastSSZ
}

func testSSVMsgJSON() []byte {
	cdJSON, _ := cd2.Encode()
	proposalData := &qbft.ProposalData{
		Data:                     cdJSON,
		RoundChangeJustification: []*qbft.SignedMessage{},
		PrepareJustification:     []*qbft.SignedMessage{},
	}

	pdJSON, _ := proposalData.Encode()

	msg := &qbft.Message{
		MsgType:    qbft.ProposalMsgType,
		Height:     1,
		Round:      2,
		Identifier: []byte{1, 2, 3, 4},
		Data:       pdJSON,
	}

	signedMsg := &qbft.SignedMessage{
		Signature: []byte{1, 2, 3, 4},
		Signers:   []OperatorID{1, 2, 3, 4},
		Message:   msg,
	}

	signedMsgJSON, _ := signedMsg.Encode()
	signedMsgJSONS2 := s2.Encode([]byte{}, signedMsgJSON)

	msgID := MessageID{}
	copy(msgID[:], msg.Identifier)

	msgToBroadcast := &SSVMessage{
		MsgType: SSVConsensusMsgType,
		MsgID:   msgID,
		Data:    signedMsgJSONS2,
	}
	msgToBroadcastJSON, _ := msgToBroadcast.Encode()
	return msgToBroadcastJSON
}

func testStructs() (Duty, phase0.AttestationData, altair.BeaconBlock, phase0.AggregateAndProof, phase0.Root, ContributionsMap) {
	dutyJSON := []byte(`{"type":0,"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`)
	attDataJSON := []byte(`{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`)
	beaconBlockJSON := []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body":{"randao_reveal":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","eth1_data":{"deposit_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":"10","block_hash":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"graffiti":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","proposer_slashings":[{"signed_header_1":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"},"signed_header_2":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x010102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}}],"attester_slashings":[{"attestation_1":{"attesting_indices":["1","2","3"],"data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"},"attestation_2":{"attesting_indices":["1","2","3"],"data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}}],"attestations":[{"aggregation_bits":"0x010203","data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}],"deposits":[{"proof":["0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f"],"data":{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000","signature":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9f"}}],"voluntary_exits":[{"message":{"epoch":"1","validator_index":"2"},"signature":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}],"sync_aggregate":{"sync_committee_bits":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","sync_committee_signature":"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f60"}}}`)
	aggregateAndProofJSON := []byte(`{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`)
	syncCommitteeContrJSON := []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x00040000000000000000000000000000","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`)

	var attDuty v1.AttesterDuty
	_ = json.Unmarshal(dutyJSON, &attDuty)
	duty := Duty{
		Type:                    4,
		PubKey:                  attDuty.PubKey,
		Slot:                    attDuty.Slot,
		ValidatorIndex:          attDuty.ValidatorIndex,
		CommitteeIndex:          attDuty.CommitteeIndex,
		CommitteeLength:         attDuty.CommitteeLength,
		CommitteesAtSlot:        attDuty.CommitteesAtSlot,
		ValidatorCommitteeIndex: attDuty.ValidatorCommitteeIndex,
	}

	var attData phase0.AttestationData
	_ = json.Unmarshal(attDataJSON, &attData)

	var beaconBlock altair.BeaconBlock
	_ = json.Unmarshal(beaconBlockJSON, &beaconBlock)

	var aggregateAndProof phase0.AggregateAndProof
	_ = json.Unmarshal(aggregateAndProofJSON, &aggregateAndProof)

	syncCommitteeBlockRoot := phase0.Root{
		0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
	}

	var syncCommitteeContr altair.SyncCommitteeContribution
	_ = json.Unmarshal(syncCommitteeContrJSON, &syncCommitteeContr)

	contMap := ContributionsMap{phase0.BLSSignature([96]byte{
		0xb1, 0x3c, 0xa7, 0x7f, 0xda, 0xb9, 0x0f, 0xce, 0xdf, 0x0c, 0xda, 0x74, 0xe9, 0xe9, 0xda, 0x1e,
		0xdb, 0xe4, 0x32, 0x91, 0x09, 0x48, 0xca, 0xad, 0xca, 0x64, 0xbb, 0xfb, 0x93, 0x34, 0x26, 0x44,
		0xac, 0xbb, 0xd3, 0xa1, 0x02, 0x4c, 0xa3, 0x9b, 0xd3, 0x50, 0x70, 0xca, 0xb3, 0xc6, 0x90, 0xd4,
		0x07, 0x43, 0x00, 0x1b, 0x44, 0x51, 0x53, 0xff, 0x97, 0x76, 0x18, 0x3c, 0xfe, 0x94, 0xec, 0x00,
		0x33, 0x90, 0xec, 0x76, 0x08, 0x4f, 0x7e, 0x20, 0x83, 0xcf, 0x3a, 0x46, 0xe1, 0xd6, 0xca, 0x1c,
		0x72, 0xb5, 0x71, 0xab, 0x58, 0x2d, 0x3d, 0x64, 0xe2, 0x69, 0x10, 0x20, 0x80, 0x85, 0x0d, 0x82,
	}): &syncCommitteeContr, phase0.BLSSignature([96]byte{
		0xb1, 0x3c, 0xa7, 0x7f, 0xda, 0xb9, 0x0f, 0xce, 0xdf, 0x0c, 0xda, 0x74, 0xe9, 0xe9, 0xda, 0x1e,
		0xdb, 0xe4, 0x32, 0x91, 0x09, 0x48, 0xca, 0xad, 0xca, 0x64, 0xbb, 0xfb, 0x93, 0x34, 0x26, 0x44,
		0xac, 0xbb, 0xd3, 0xa1, 0x02, 0x4c, 0xa3, 0x9b, 0xd3, 0x50, 0x70, 0xca, 0xb3, 0xc6, 0x90, 0xd4,
		0x07, 0x43, 0x00, 0x1b, 0x44, 0x51, 0x53, 0xff, 0x97, 0x76, 0x18, 0x3c, 0xfe, 0x94, 0xec, 0x00,
		0x33, 0x90, 0xec, 0x76, 0x08, 0x4f, 0x7e, 0x20, 0x83, 0xcf, 0x3a, 0x46, 0xe1, 0xd6, 0xca, 0x1c,
		0x72, 0xb5, 0x71, 0xbb, 0x58, 0x2d, 0x3d, 0x64, 0xe2, 0x69, 0x10, 0x20, 0x80, 0x85, 0x0d, 0x82,
	}): &syncCommitteeContr, phase0.BLSSignature([96]byte{
		0xb1, 0x3c, 0xa7, 0x7f, 0xda, 0xb9, 0x0f, 0xce, 0xdf, 0x0c, 0xda, 0x74, 0xe9, 0xe9, 0xda, 0x1e,
		0xdb, 0xe4, 0x32, 0x91, 0x09, 0x48, 0xca, 0xad, 0xca, 0x64, 0xbb, 0xfb, 0x93, 0x34, 0x26, 0x44,
		0xac, 0xbb, 0xd3, 0xa1, 0x02, 0x4c, 0xa3, 0x9b, 0xd3, 0x50, 0x70, 0xca, 0xb3, 0xc6, 0x90, 0xd4,
		0x07, 0x43, 0x00, 0x1b, 0x44, 0x51, 0x53, 0xff, 0x97, 0x76, 0x18, 0x3c, 0xfe, 0x94, 0xec, 0x00,
		0x33, 0x90, 0xec, 0x76, 0x08, 0x4f, 0x7e, 0x20, 0x83, 0xcf, 0x3a, 0x46, 0xe1, 0xd6, 0xca, 0x1c,
		0x72, 0xb5, 0x71, 0xbb, 0x58, 0x22, 0x3d, 0x64, 0xe2, 0x69, 0x10, 0x20, 0x80, 0x85, 0x0d, 0x82,
	}): &syncCommitteeContr, phase0.BLSSignature([96]byte{
		0xb1, 0x3c, 0xa7, 0x7f, 0xda, 0xb9, 0x0f, 0xce, 0xdf, 0x0c, 0xda, 0x74, 0xe9, 0xe9, 0xda, 0x1e,
		0xdb, 0xe4, 0x32, 0x91, 0x09, 0x48, 0xca, 0xad, 0xca, 0x64, 0xbb, 0xfb, 0x93, 0x34, 0x26, 0x44,
		0xac, 0xbb, 0xd3, 0xa1, 0x02, 0x4c, 0xa3, 0x9b, 0xd3, 0x50, 0x70, 0xca, 0xb3, 0xc6, 0x90, 0xd4,
		0x07, 0x43, 0x00, 0x1b, 0x44, 0x51, 0x53, 0xff, 0x97, 0x76, 0x18, 0x3c, 0xfe, 0x94, 0xec, 0x00,
		0x33, 0x90, 0xec, 0x76, 0x08, 0x4f, 0x7e, 0x20, 0x83, 0xcf, 0x3a, 0x46, 0xe1, 0xd6, 0xca, 0x1c,
		0x72, 0xb5, 0x71, 0xbb, 0x58, 0x23, 0x3d, 0x64, 0xe2, 0x69, 0x10, 0x20, 0x80, 0x85, 0x0d, 0x99,
	}): &syncCommitteeContr}

	return duty, attData, beaconBlock, aggregateAndProof, syncCommitteeBlockRoot, contMap
}

func BenchmarkEncodeSSZ(b *testing.B) {
	var s []byte
	//duty, attData, _, _, _, _ := testStructs()
	//cd2 := ConsensusData{
	//	Duty:            &duty,
	//	AttestationData: &attData,
	//}
	for i := 0; i < b.N; i++ {
		s = testSSVMsgSSZ()
		//s, _ = cd2.MarshalSSZ()
	}
	result = s
}

//func BenchmarkSSZCompressSnappy(b *testing.B) {
//	var s []byte
//	//duty, attData, _, _, _, _ := testStructs()
//	//cd2 := ConsensusData{
//	//	Duty:            &duty,
//	//	AttestationData: &attData,
//	//}
//	for i := 0; i < b.N; i++ {
//		s, _ = cd2.MarshalSSZ()
//		s = snappy.Encode([]byte{}, s)
//		//s, _ = snappy.Encode(nil, s)
//	}
//	result = s
//}

//func BenchmarkSSZCompressS2Snappy(b *testing.B) {
//	var s []byte
//	//duty, attData, _, _, _, _ := testStructs()
//	//cd2 := ConsensusData{
//	//	Duty:            &duty,
//	//	AttestationData: &attData,
//	//}
//	for i := 0; i < b.N; i++ {
//		s, _ = cd2.MarshalSSZ()
//		s = s2.EncodeSnappy([]byte{}, s)
//	}
//	result = s
//}

//func BenchmarkSSZCompressS2(b *testing.B) {
//	var s []byte
//	//duty, attData, _, _, _, _ := testStructs()
//	//cd2 := ConsensusData{
//	//	Duty:            &duty,
//	//	AttestationData: &attData,
//	//}
//	for i := 0; i < b.N; i++ {
//		s, _ = cd2.MarshalSSZ()
//		s = s2.Encode([]byte{}, s)
//	}
//	result = s
//}

func BenchmarkEncodeJSON(b *testing.B) {
	var s []byte
	for i := 0; i < b.N; i++ {
		s = testSSVMsgJSON()
	}
	result = s
}

//func BenchmarkJSONCompressSnappy(b *testing.B) {
//	var s []byte
//	//duty, attData, _, _, _, _ := testStructs()
//	//cd2 := ConsensusData{
//	//	Duty:            &duty,
//	//	AttestationData: &attData,
//	//}
//	for i := 0; i < b.N; i++ {
//		s, _ = cd2.Encode()
//		s = snappy.Encode([]byte{}, s)
//		//s, _ = snappy.Encode(nil, s)
//	}
//	result = s
//}

//func BenchmarkJSONCompressS2Snappy(b *testing.B) {
//	var s []byte
//	//duty, attData, _, _, _, _ := testStructs()
//	//cd2 := ConsensusData{
//	//	Duty:            &duty,
//	//	AttestationData: &attData,
//	//}
//	for i := 0; i < b.N; i++ {
//		s, _ = cd2.Encode()
//		s = s2.EncodeSnappy([]byte{}, s)
//	}
//	result = s
//}

//func BenchmarkJSONCompressS2(b *testing.B) {
//	var s []byte
//	//duty, attData, _, _, _, _ := testStructs()
//	//cd2 := ConsensusData{
//	//	Duty:            &duty,
//	//	AttestationData: &attData,
//	//}
//	for i := 0; i < b.N; i++ {
//		s, _ = cd2.Encode()
//		s = s2.Encode([]byte{}, s)
//	}
//	result = s
//}

func TestEncoding(t *testing.T) {
	t.Skip()
	dutyJSON := []byte(`{"type":0,"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","slot":"1","validator_index":"2","committee_index":"3","committee_length":"128","committees_at_slot":"4","validator_committee_index":"61"}`)
	attDataJSON := []byte(`{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}}`)
	beaconBlockJSON := []byte(`{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body":{"randao_reveal":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","eth1_data":{"deposit_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","deposit_count":"10","block_hash":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"graffiti":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","proposer_slashings":[{"signed_header_1":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"},"signed_header_2":{"message":{"slot":"1","proposer_index":"2","parent_root":"0x010102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","state_root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","body_root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}}],"attester_slashings":[{"attestation_1":{"attesting_indices":["1","2","3"],"data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"},"attestation_2":{"attesting_indices":["1","2","3"],"data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}}],"attestations":[{"aggregation_bits":"0x010203","data":{"slot":"100","index":"1","beacon_block_root":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","source":{"epoch":"1","root":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f"},"target":{"epoch":"2","root":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}},"signature":"0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9fa0a1a2a3a4a5a6a7a8a9aaabacadaeafb0b1b2b3b4b5b6b7b8b9babbbcbdbebf"}],"deposits":[{"proof":["0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f","0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f","0x606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f"],"data":{"pubkey":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f","withdrawal_credentials":"0x202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","amount":"32000000000","signature":"0x404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f606162636465666768696a6b6c6d6e6f707172737475767778797a7b7c7d7e7f808182838485868788898a8b8c8d8e8f909192939495969798999a9b9c9d9e9f"}}],"voluntary_exits":[{"message":{"epoch":"1","validator_index":"2"},"signature":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f"}],"sync_aggregate":{"sync_committee_bits":"0x000102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f","sync_committee_signature":"0x0102030405060708090a0b0c0d0e0f101112131415161718191a1b1c1d1e1f202122232425262728292a2b2c2d2e2f303132333435363738393a3b3c3d3e3f404142434445464748494a4b4c4d4e4f505152535455565758595a5b5c5d5e5f60"}}}`)
	aggregateAndProofJSON := []byte(`{"aggregator_index":"402","aggregate":{"aggregation_bits":"0xffffffff01","data":{"slot":"66","index":"0","beacon_block_root":"0x737b2949b471552a7f95f772e289ae6d74bd8e527120d9993095fd34ed89e100","source":{"epoch":"0","root":"0x0000000000000000000000000000000000000000000000000000000000000000"},"target":{"epoch":"2","root":"0x674d7e0ce7a28ba0d71ecef8d44621e8f4ed206e9116dc647fafd7f32f61f440"}},"signature":"0x8a75731b877a4be72ddc81ae5318eaa9863fef2297b58a4f01a447bd1fff10d48bb79e62d280557c472af5d457032e0112db17f99b2e925ce2c89dd839e5bd8e5e95b2f5253bb80087753555c69b116162c334f5a142e38ff6a66ef579c9a70d"},"selection_proof":"0x8b5f33a895612754103fbaaed74b408e89b948c69740d722b56207c272e001b2ddd445931e40a2938c84afab86c2606f0c1a93a0aaf4962c91d3ddf309de8ef0dbd68f590573e53e5ff7114e9625fae2cfee9e7eb991ad929d351c7701581d9c"}`)
	syncCommitteeContrJSON := []byte(`{"slot":"1","beacon_block_root":"0xbacd20f09da907734434f052bd4c9503aa16bab1960e89ea20610d08d064481c","subcommittee_index":"3","aggregation_bits":"0x00040000000000000000000000000000","signature":"0xb591bd4ca7d745b6e027879645d7c014fecb8c58631af070f7607acc0c1c948a5102a33267f0e4ba41a85b254b07df91185274375b2e6436e37e81d2fd46cb3751f5a6c86efb7499c1796c0c17e122a54ac067bb0f5ff41f3241659cceb0c21c"}`)

	var duty v1.AttesterDuty
	err := json.Unmarshal(dutyJSON, &duty)
	require.NoError(t, err)

	var attData phase0.AttestationData
	err = json.Unmarshal(attDataJSON, &attData)
	require.NoError(t, err)

	var beaconBlock altair.BeaconBlock
	err = json.Unmarshal(beaconBlockJSON, &beaconBlock)
	require.NoError(t, err)

	var aggregateAndProof phase0.AggregateAndProof
	err = json.Unmarshal(aggregateAndProofJSON, &aggregateAndProof)
	require.NoError(t, err)

	var syncCommitteeContr altair.SyncCommitteeContribution
	err = json.Unmarshal(syncCommitteeContrJSON, &syncCommitteeContr)
	require.NoError(t, err)

	cd := ConsensusData{
		Duty: &Duty{
			Type:                    4,
			PubKey:                  duty.PubKey,
			Slot:                    duty.Slot,
			ValidatorIndex:          duty.ValidatorIndex,
			CommitteeIndex:          duty.CommitteeIndex,
			CommitteeLength:         duty.CommitteeLength,
			CommitteesAtSlot:        duty.CommitteesAtSlot,
			ValidatorCommitteeIndex: duty.ValidatorCommitteeIndex,
		},
		//AttestationData: &attData,
		//BlockData: &beaconBlock,
		//AggregateAndProof: &aggregateAndProof,
		//SyncCommitteeBlockRoot: phase0.Root{
		//	0x00, 0x01, 0x02, 0x03, 0x04, 0x05, 0x06, 0x07, 0x08, 0x09, 0x0a, 0x0b, 0x0c, 0x0d, 0x0e, 0x0f,
		//	0x10, 0x11, 0x12, 0x13, 0x14, 0x15, 0x16, 0x17, 0x18, 0x19, 0x1a, 0x1b, 0x1c, 0x1d, 0x1e, 0x1f,
		//},
		SyncCommitteeContribution: ContributionsMap{phase0.BLSSignature([96]byte{
			0xb1, 0x3c, 0xa7, 0x7f, 0xda, 0xb9, 0x0f, 0xce, 0xdf, 0x0c, 0xda, 0x74, 0xe9, 0xe9, 0xda, 0x1e,
			0xdb, 0xe4, 0x32, 0x91, 0x09, 0x48, 0xca, 0xad, 0xca, 0x64, 0xbb, 0xfb, 0x93, 0x34, 0x26, 0x44,
			0xac, 0xbb, 0xd3, 0xa1, 0x02, 0x4c, 0xa3, 0x9b, 0xd3, 0x50, 0x70, 0xca, 0xb3, 0xc6, 0x90, 0xd4,
			0x07, 0x43, 0x00, 0x1b, 0x44, 0x51, 0x53, 0xff, 0x97, 0x76, 0x18, 0x3c, 0xfe, 0x94, 0xec, 0x00,
			0x33, 0x90, 0xec, 0x76, 0x08, 0x4f, 0x7e, 0x20, 0x83, 0xcf, 0x3a, 0x46, 0xe1, 0xd6, 0xca, 0x1c,
			0x72, 0xb5, 0x71, 0xab, 0x58, 0x2d, 0x3d, 0x64, 0xe2, 0x69, 0x10, 0x20, 0x80, 0x85, 0x0d, 0x82,
		}): &syncCommitteeContr, phase0.BLSSignature([96]byte{
			0xb1, 0x3c, 0xa7, 0x7f, 0xda, 0xb9, 0x0f, 0xce, 0xdf, 0x0c, 0xda, 0x74, 0xe9, 0xe9, 0xda, 0x1e,
			0xdb, 0xe4, 0x32, 0x91, 0x09, 0x48, 0xca, 0xad, 0xca, 0x64, 0xbb, 0xfb, 0x93, 0x34, 0x26, 0x44,
			0xac, 0xbb, 0xd3, 0xa1, 0x02, 0x4c, 0xa3, 0x9b, 0xd3, 0x50, 0x70, 0xca, 0xb3, 0xc6, 0x90, 0xd4,
			0x07, 0x43, 0x00, 0x1b, 0x44, 0x51, 0x53, 0xff, 0x97, 0x76, 0x18, 0x3c, 0xfe, 0x94, 0xec, 0x00,
			0x33, 0x90, 0xec, 0x76, 0x08, 0x4f, 0x7e, 0x20, 0x83, 0xcf, 0x3a, 0x46, 0xe1, 0xd6, 0xca, 0x1c,
			0x72, 0xb5, 0x71, 0xbb, 0x58, 0x2d, 0x3d, 0x64, 0xe2, 0x69, 0x10, 0x20, 0x80, 0x85, 0x0d, 0x82,
		}): &syncCommitteeContr, phase0.BLSSignature([96]byte{
			0xb1, 0x3c, 0xa7, 0x7f, 0xda, 0xb9, 0x0f, 0xce, 0xdf, 0x0c, 0xda, 0x74, 0xe9, 0xe9, 0xda, 0x1e,
			0xdb, 0xe4, 0x32, 0x91, 0x09, 0x48, 0xca, 0xad, 0xca, 0x64, 0xbb, 0xfb, 0x93, 0x34, 0x26, 0x44,
			0xac, 0xbb, 0xd3, 0xa1, 0x02, 0x4c, 0xa3, 0x9b, 0xd3, 0x50, 0x70, 0xca, 0xb3, 0xc6, 0x90, 0xd4,
			0x07, 0x43, 0x00, 0x1b, 0x44, 0x51, 0x53, 0xff, 0x97, 0x76, 0x18, 0x3c, 0xfe, 0x94, 0xec, 0x00,
			0x33, 0x90, 0xec, 0x76, 0x08, 0x4f, 0x7e, 0x20, 0x83, 0xcf, 0x3a, 0x46, 0xe1, 0xd6, 0xca, 0x1c,
			0x72, 0xb5, 0x71, 0xbb, 0x58, 0x22, 0x3d, 0x64, 0xe2, 0x69, 0x10, 0x20, 0x80, 0x85, 0x0d, 0x82,
		}): &syncCommitteeContr, phase0.BLSSignature([96]byte{
			0xb1, 0x3c, 0xa7, 0x7f, 0xda, 0xb9, 0x0f, 0xce, 0xdf, 0x0c, 0xda, 0x74, 0xe9, 0xe9, 0xda, 0x1e,
			0xdb, 0xe4, 0x32, 0x91, 0x09, 0x48, 0xca, 0xad, 0xca, 0x64, 0xbb, 0xfb, 0x93, 0x34, 0x26, 0x44,
			0xac, 0xbb, 0xd3, 0xa1, 0x02, 0x4c, 0xa3, 0x9b, 0xd3, 0x50, 0x70, 0xca, 0xb3, 0xc6, 0x90, 0xd4,
			0x07, 0x43, 0x00, 0x1b, 0x44, 0x51, 0x53, 0xff, 0x97, 0x76, 0x18, 0x3c, 0xfe, 0x94, 0xec, 0x00,
			0x33, 0x90, 0xec, 0x76, 0x08, 0x4f, 0x7e, 0x20, 0x83, 0xcf, 0x3a, 0x46, 0xe1, 0xd6, 0xca, 0x1c,
			0x72, 0xb5, 0x71, 0xbb, 0x58, 0x23, 0x3d, 0x64, 0xe2, 0x69, 0x10, 0x20, 0x80, 0x85, 0x0d, 0x99,
		}): &syncCommitteeContr},
	}

	cdSSZ, err := cd.MarshalSSZ()
	require.NoError(t, err)

	//cdSSZ, err := cd.MarshalSSZ()
	//require.NoError(t, err)

	var unmarshalCd ConsensusData
	err = unmarshalCd.UnmarshalSSZ(cdSSZ)
	require.NoError(t, err)
	//litter.Dump(unmarshalCd)

	require.Equal(t, cd, unmarshalCd)

	sszSnappyByts := snappy.Encode([]byte{}, cdSSZ)
	sszS2Byts := s2.Encode([]byte{}, cdSSZ)
	jsonByts, _ := cd.Encode()
	jsonSnappyByts := snappy.Encode([]byte{}, jsonByts)
	jsonS2Byts := s2.Encode([]byte{}, jsonByts)

	fmt.Printf("ssz:         %d\n", len(cdSSZ))
	fmt.Printf("ssz:         %d\n", len(cdSSZ))
	fmt.Printf("ssz snappy:  %d\n", len(sszSnappyByts))
	fmt.Printf("ssz s2:      %d\n", len(sszS2Byts))
	fmt.Printf("json:        %d\n", len(jsonByts))
	fmt.Printf("json snappy: %d\n", len(jsonSnappyByts))
	fmt.Printf("json s2:     %d\n", len(jsonS2Byts))
}
