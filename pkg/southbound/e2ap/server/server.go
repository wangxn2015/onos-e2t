// SPDX-FileCopyrightText: 2020-present Open Networking Foundation <info@opennetworking.org>
//
// SPDX-License-Identifier: Apache-2.0

package server

import (
	"context"
	"encoding/hex"
	"fmt"
	ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
	"strconv"
	"time"

	"github.com/wangxn2015/onos-e2t/pkg/southbound/e2ap/stream"
	"github.com/wangxn2015/onos-lib-go/pkg/logging"

	e2ap_ies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"

	prototypes "github.com/gogo/protobuf/types"

	"github.com/wangxn2015/onos-e2t/pkg/store/rnib"

	"github.com/onosproject/onos-e2t/api/e2ap/v2"
	e2apies "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-ies"

	e2smtypes "github.com/onosproject/onos-api/go/onos/e2t/e2sm"

	topoapi "github.com/onosproject/onos-api/go/onos/topo"
	e2appducontents "github.com/onosproject/onos-e2t/api/e2ap/v2/e2ap-pdu-contents"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/pdubuilder"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/pdudecoder"
	"github.com/onosproject/onos-e2t/pkg/southbound/e2ap/types"
	"github.com/wangxn2015/onos-e2t/pkg/modelregistry"
	e2 "github.com/wangxn2015/onos-e2t/pkg/protocols/e2ap"
	//"github.com/onosproject/ran-simulator/pkg/utils/e2sm/kpm2/id/cellglobalid"
	"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
)

var log = logging.GetLogger()

// TODO: Change the RIC ID to something appropriate
var ricID = types.RicIdentifier{
	RicIdentifierValue: []byte{0xDE, 0xBC, 0xA0},
	RicIdentifierLen:   20,
}

func NewE2Server(e2apConns E2APConnManager,
	mgmtConns MgmtConnManager,
	streams stream.Manager,
	modelRegistry modelregistry.ModelRegistry, rnib rnib.Store) *E2Server {
	return &E2Server{
		server:    e2.NewServer(),
		e2apConns: e2apConns,
		mgmtConns: mgmtConns,

		streams:       streams,
		modelRegistry: modelRegistry,
		rnib:          rnib,
	}
}

type E2Server struct {
	server        *e2.Server
	e2apConns     E2APConnManager
	mgmtConns     MgmtConnManager
	streams       stream.Manager
	modelRegistry modelregistry.ModelRegistry
	rnib          rnib.Store
}

func (s *E2Server) Serve() error {
	return s.server.Serve(func(conn e2.ServerConn) e2.ServerInterface {
		return &E2APServer{
			serverConn:    conn,
			e2apConns:     s.e2apConns,
			mgmtConns:     s.mgmtConns,
			streams:       s.streams,
			modelRegistry: s.modelRegistry,
			rnib:          s.rnib,
		}
	})
}

func (s *E2Server) Stop() error {
	return s.server.Stop()
}

type E2APServer struct {
	e2apConns     E2APConnManager
	mgmtConns     MgmtConnManager
	streams       stream.Manager
	serverConn    e2.ServerConn
	e2apConn      *E2APConn
	modelRegistry modelregistry.ModelRegistry
	rnib          rnib.Store
}

// uint24ToUint32 converts uint24 uint32
func uint24ToUint32(val []byte) uint32 {
	r := uint32(0)
	for i := uint32(0); i < 3; i++ {
		r |= uint32(val[i]) << (8 * i)
	}
	return r
}

func (e *E2APServer) E2Setup(ctx context.Context, request *e2appducontents.E2SetupRequest) (*e2appducontents.E2SetupResponse, *e2appducontents.E2SetupFailure, error) {
	log.Warnf("Received E2 setup request: %+v", request)
	transID, nodeIdentity, ranFuncs, _, err := pdudecoder.DecodeE2SetupRequest(request)
	if err != nil {
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}

		var trID int32
		for _, v := range request.GetProtocolIes() {
			if v.Id == int32(v2.ProtocolIeIDTransactionID) {
				trID = v.GetValue().GetTrId().GetValue()
			}
		}

		failure := &e2appducontents.E2SetupFailure{
			ProtocolIes: make([]*e2appducontents.E2SetupFailureIes, 0),
		}
		failure.SetErrorCause(cause).SetTransactionID(trID)

		return nil, failure, err
	}

	rawPlmnid := []byte{nodeIdentity.Plmn[0], nodeIdentity.Plmn[1], nodeIdentity.Plmn[2]}
	plmnID := fmt.Sprintf("%x", uint24ToUint32(rawPlmnid))

	var e2Cells []*topoapi.E2Cell
	serviceModels := make(map[string]*topoapi.ServiceModelInfo)
	rfAccepted := make(types.RanFunctionRevisions)
	rfRejected := make(types.RanFunctionCauses)
	plugins := e.modelRegistry.GetPlugins()
	//第一层，遍历e2t中的sm插件
	for smOid, sm := range plugins {
		// 对每一个sm插件
		var ranFunctions []*prototypes.Any
		var ranFunctionIDs []uint32
		//第二层， 遍历 本node在e2setup request中的 RAN func.  变量ranFuncs理解为list,每一个元素为 ran func item的简化修改版本
		for ranFunctionID, ranFunc := range *ranFuncs {
			oid := e2smtypes.OID(ranFunc.OID)
			if smOid == oid { //e.g. kpm plugin meets with E2Node ran func
				serviceModels[string(smOid)] = &topoapi.ServiceModelInfo{ //新建一个
					OID:          string(smOid),
					RanFunctions: ranFunctions, //
				}
				//存储1个id   如 kpm 的 ranFunctionID 为 4
				ranFunctionIDs = append(ranFunctionIDs, uint32(ranFunctionID))
				// e.g. 将kpm sm 转化为 modelregistry.E2Setup接口
				if setup, ok := sm.(modelregistry.E2Setup); ok {
					//生成onSetupRequest，作为插件OnSetup函数的参数
					onSetupRequest := &e2smtypes.OnSetupRequest{
						ServiceModels:          serviceModels,
						E2Cells:                &e2Cells,
						RANFunctionDescription: ranFunc.Description,
					}
					err := setup.OnSetup(onSetupRequest)

					//----------------------------
					//!!! -------added by wxn to compensate BAICELLS ran e2ap version problem. 2022.10.14
					//17号周一测试 此代码欲解决问题baicell ran 不建立cell的问题
					//!!! -------ran's implementation skips kpm nodelist since it's only a kmpv203 instead of kpmv2.0
					//----------------------------
					//
					if len(*onSetupRequest.E2Cells) == 0 {
						log.Warnf("wxn----> length of onSetupRequest.E2Cells is zero")
						cellNcgi := uint64(87893173159116801)
						nci := cellNcgi & 0xfffffffff
						ncibs := &asn1.BitString{
							Value: Uint64ToBitString(nci, 36),
							Len:   36,
						}

						plmnIDx := ransimtypes.NewUint24(uint32(1279014))

						// move the function into utils.go under server dir
						cellGlobalID := NewGlobalNRCGIID(WithPlmnID(plmnIDx), WithNRCellID(ncibs))

						if err != nil {
							log.Warn(err)
						}

						// 生成cellObject.CellObjectID
						cellObject := &topoapi.E2Cell{
							CellObjectID: strconv.FormatUint(uint64(cellNcgi), 16),
							CellGlobalID: &topoapi.CellGlobalID{
								Value: fmt.Sprintf("%x", bitStringToUint64(cellGlobalID.nrCellID.Value, int(cellGlobalID.nrCellID.Len))),
								Type:  topoapi.CellGlobalIDType_NRCGI,
							},
						}
						//添加到内部
						*onSetupRequest.E2Cells = append(*onSetupRequest.E2Cells, cellObject)
					}

					//------------------------end ------------------------------
					//------------------------------------------------------

					if err != nil {
						log.Warn(err)
						log.Warnf("Length of RAN function Description Bytes is: %d", len(onSetupRequest.RANFunctionDescription))
						log.Warnf("RAN Function Description Bytes in hex format: %v", hex.Dump(onSetupRequest.RANFunctionDescription))
					}

				}
				rfAccepted[ranFunctionID] = ranFunc.Revision
				serviceModels[string(smOid)].RanFunctionIDs = ranFunctionIDs
			}
		}
	}

	mgmtConn := NewMgmtConn(createE2NodeURI(nodeIdentity), plmnID, nodeIdentity, e.serverConn, serviceModels, e2Cells, time.Now())

	// Create an E2 setup response
	e2ncID3 := pdubuilder.CreateE2NodeComponentIDS1("S1-component")
	e2nccaal := make([]*types.E2NodeComponentConfigAdditionAckItem, 0)
	ie1 := types.E2NodeComponentConfigAdditionAckItem{
		E2NodeComponentConfigurationAck: e2ap_ies.E2NodeComponentConfigurationAck{
			UpdateOutcome: e2ap_ies.UpdateOutcome_UPDATE_OUTCOME_SUCCESS,
		},
		E2NodeComponentID:   e2ncID3,
		E2NodeComponentType: e2ap_ies.E2NodeComponentInterfaceType_E2NODE_COMPONENT_INTERFACE_TYPE_S1,
	}
	e2nccaal = append(e2nccaal, &ie1)
	response, err := pdubuilder.NewE2SetupResponse(*transID, nodeIdentity.Plmn, ricID, e2nccaal)
	if err != nil {
		cause := &e2apies.Cause{
			Cause: &e2apies.Cause_RicRequest{
				RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
			},
		}

		var trID int32
		for _, v := range request.GetProtocolIes() {
			if v.Id == int32(v2.ProtocolIeIDTransactionID) {
				trID = v.GetValue().GetTrId().GetValue()
			}
		}

		failure := &e2appducontents.E2SetupFailure{
			ProtocolIes: make([]*e2appducontents.E2SetupFailureIes, 0),
		}
		failure.SetErrorCause(cause).SetTransactionID(trID)

		return nil, failure, err
	}

	if len(rfAccepted) > 0 {
		response.SetRanFunctionAccepted(rfAccepted)
	}
	if len(rfRejected) > 0 {
		response.SetRanFunctionRejected(rfRejected)
	}
	log.Warnf("Sending E2 setup response %+v", response)
	e.mgmtConns.open(mgmtConn)
	return response, nil, nil
}

func (e *E2APServer) RICIndication(ctx context.Context, request *e2appducontents.Ricindication) error {
	return e.e2apConn.ricIndication(ctx, request)
}

func (e *E2APServer) E2ConfigurationUpdate(ctx context.Context, request *e2appducontents.E2NodeConfigurationUpdate) (response *e2appducontents.E2NodeConfigurationUpdateAcknowledge, failure *e2appducontents.E2NodeConfigurationUpdateFailure, err error) {
	log.Warnf("Received E2 node configuration update request: %+v", request)

	var nodeIdentity *e2apies.GlobalE2NodeId
	e2nccual := make([]*types.E2NodeComponentConfigUpdateItem, 0)
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDGlobalE2nodeID) {
			nodeIdentity = v.GetValue().GetGe2NId()
		}
		if v.Id == int32(v2.ProtocolIeIDE2nodeComponentConfigUpdate) {
			list := v.GetValue().GetE2Nccul().GetValue()
			for _, ie := range list {
				e2nccuai := types.E2NodeComponentConfigUpdateItem{}
				e2nccuai.E2NodeComponentType = ie.GetValue().GetE2Nccui().GetE2NodeComponentInterfaceType()
				e2nccuai.E2NodeComponentID = ie.GetValue().GetE2Nccui().GetE2NodeComponentId()
				e2nccuai.E2NodeComponentConfiguration = *ie.GetValue().GetE2Nccui().GetE2NodeComponentConfiguration()

				e2nccual = append(e2nccual, &e2nccuai)
			}
		}
	}

	if nodeIdentity != nil {
		// --------------wxn TO DO: simplify E2 config update processing
		nodeID, err := pdudecoder.ExtractE2NodeIdentity(nodeIdentity, e2nccual)
		if err != nil {
			cause := &e2apies.Cause{
				Cause: &e2apies.Cause_RicRequest{
					RicRequest: e2apies.CauseRicrequest_CAUSE_RICREQUEST_UNSPECIFIED,
				},
			}

			var trID int32
			for _, v := range request.GetProtocolIes() {
				if v.Id == int32(v2.ProtocolIeIDTransactionID) {
					trID = v.GetValue().GetTrId().GetValue()
					break
				}
			}

			failure := &e2appducontents.E2NodeConfigurationUpdateFailure{
				ProtocolIes: make([]*e2appducontents.E2NodeConfigurationUpdateFailureIes, 0),
			}
			failure.SetCause(cause).SetTransactionID(trID)

			return nil, failure, nil
		}

		// Creates a new E2AP data connection ----wxn : note here
		e.e2apConn = NewE2APConn(createE2NodeURI(nodeID), e.serverConn, e.streams, e.rnib)
	}

	var trID int32
	for _, v := range request.GetProtocolIes() {
		if v.Id == int32(v2.ProtocolIeIDTransactionID) {
			trID = v.GetValue().GetTrId().GetValue()
			break
		}
	}

	e2ncua := &e2appducontents.E2NodeConfigurationUpdateAcknowledge{
		ProtocolIes: make([]*e2appducontents.E2NodeConfigurationUpdateAcknowledgeIes, 0),
	}
	e2ncua.SetTransactionID(trID)
	log.Warnf("Composed E2nodeConfigurationUpdateMessage is\n%v", e2ncua)
	log.Warnf("Sending config update ack to e2 node: %s", e.e2apConn.E2NodeID)
	//--------------wxn : note here
	e.e2apConns.open(e.e2apConn)
	return e2ncua, nil, nil
}

func bitStringToUint64(bitString []byte, bitCount int) uint64 {
	var result uint64
	for i, b := range bitString {
		result += uint64(b) << ((len(bitString) - i - 1) * 8)
	}
	if bitCount%8 != 0 {
		return result >> (8 - bitCount%8)
	}
	return result
}

//-----------------------
//package cellglobalid
//
//import (
//ransimtypes "github.com/onosproject/onos-api/go/onos/ransim/types"
//e2smkpmv2 "github.com/onosproject/onos-e2-sm/servicemodels/e2sm_kpm_v2_go/v2/e2sm-kpm-v2-go"
//"github.com/onosproject/onos-lib-go/api/asn1/v1/asn1"
//)

// GlobalNRCGIID cell global NRCGI ID
type GlobalNRCGIID struct {
	plmnID   *ransimtypes.Uint24
	nrCellID *asn1.BitString
}

// NewGlobalNRCGIID creates new global NRCGI ID
func NewGlobalNRCGIID(options ...func(*GlobalNRCGIID)) *GlobalNRCGIID {
	nrcgiid := &GlobalNRCGIID{}
	for _, option := range options {
		option(nrcgiid)
	}

	return nrcgiid
}

//
// WithPlmnID sets plmn ID
//                                 WithPlmnID函数 返回一个函数：  func(nrcgiid *GlobalNRCGIID）
func WithPlmnID(plmnID *ransimtypes.Uint24) func(nrcgiid *GlobalNRCGIID) {
	return func(nrcgid *GlobalNRCGIID) {
		nrcgid.plmnID = plmnID

	}
}

// WithNRCellID sets NRCellID
func WithNRCellID(nrCellID *asn1.BitString) func(nrcgiid *GlobalNRCGIID) {
	return func(nrcgid *GlobalNRCGIID) {
		nrcgid.nrCellID = nrCellID
	}
}

//// Build builds a global NRCGI ID
//func (gNRCGIID *GlobalNRCGIID) Build() (*e2smkpmv2.CellGlobalId, error) {
//	return &e2smkpmv2.CellGlobalId{
//		CellGlobalId: &e2smkpmv2.CellGlobalId_NrCgi{
//			NrCgi: &e2smkpmv2.Nrcgi{
//				PLmnIdentity: &e2smkpmv2.PlmnIdentity{
//					Value: gNRCGIID.plmnID.ToBytes(),
//				},
//				NRcellIdentity: &e2smkpmv2.NrcellIdentity{
//					Value: gNRCGIID.nrCellID,
//				},
//			},
//		},
//	}, nil
//}

// Uint64ToBitString converts uint64 to a bit string byte array
func Uint64ToBitString(value uint64, bitCount int) []byte {
	result := make([]byte, bitCount/8+1)
	if bitCount%8 > 0 {
		value = value << (8 - bitCount%8)
	}

	for i := 0; i <= (bitCount / 8); i++ {
		result[i] = byte(value >> (((bitCount / 8) - i) * 8) & 0xFF)
	}

	return result
}
