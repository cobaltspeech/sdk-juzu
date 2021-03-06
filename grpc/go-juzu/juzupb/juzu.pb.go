// Code generated by protoc-gen-go. DO NOT EDIT.
// source: juzu.proto

package juzupb

import (
	context "context"
	fmt "fmt"
	proto "github.com/golang/protobuf/proto"
	duration "github.com/golang/protobuf/ptypes/duration"
	empty "github.com/golang/protobuf/ptypes/empty"
	_ "google.golang.org/genproto/googleapis/api/annotations"
	grpc "google.golang.org/grpc"
	math "math"
)

// Reference imports to suppress errors if they are not otherwise used.
var _ = proto.Marshal
var _ = fmt.Errorf
var _ = math.Inf

// This is a compile-time assertion to ensure that this generated file
// is compatible with the proto package it is being compiled against.
// A compilation error at this line likely means your copy of the
// proto package needs to be updated.
const _ = proto.ProtoPackageIsVersion3 // please upgrade the proto package

// The encoding of the audio data to be sent for recognition.
//
// For best results, the audio source should be captured and transmitted using
// the RAW_LINEAR16 encoding.
type DiarizationConfig_Encoding int32

const (
	// Raw (headerless) Uncompressed 16-bit signed little endian samples (linear
	// PCM), single channel, sampled at the rate expected by the chosen `Model`.
	DiarizationConfig_RAW_LINEAR16 DiarizationConfig_Encoding = 0
	// WAV (data with RIFF headers), with data sampled at a rate equal to or
	// higher than the sample rate expected by the chosen Model.
	DiarizationConfig_WAV DiarizationConfig_Encoding = 1
	// FLAC data, sampled at a rate equal to or higher than the sample rate
	// expected by the chosen Model.
	DiarizationConfig_FLAC DiarizationConfig_Encoding = 2
)

var DiarizationConfig_Encoding_name = map[int32]string{
	0: "RAW_LINEAR16",
	1: "WAV",
	2: "FLAC",
}

var DiarizationConfig_Encoding_value = map[string]int32{
	"RAW_LINEAR16": 0,
	"WAV":          1,
	"FLAC":         2,
}

func (x DiarizationConfig_Encoding) String() string {
	return proto.EnumName(DiarizationConfig_Encoding_name, int32(x))
}

func (DiarizationConfig_Encoding) EnumDescriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{6, 0}
}

// The top-level message sent by the client for the `StreamingDiarize`
// request.  Multiple `StreamingDiarizeRequest` messages are sent. The first
// message must contain a `DiarizationConfig` message only, and all subsequent
// messages must contain `DiarizationAudio ` only.  All `DiarizationAudio `
// messages must contain non-empty audio.  If audio content is empty, the server
// may interpret it as end of stream and stop accepting any further messages.
type StreamingDiarizeRequest struct {
	// Types that are valid to be assigned to Request:
	//	*StreamingDiarizeRequest_Config
	//	*StreamingDiarizeRequest_Audio
	Request              isStreamingDiarizeRequest_Request `protobuf_oneof:"request"`
	XXX_NoUnkeyedLiteral struct{}                          `json:"-"`
	XXX_unrecognized     []byte                            `json:"-"`
	XXX_sizecache        int32                             `json:"-"`
}

func (m *StreamingDiarizeRequest) Reset()         { *m = StreamingDiarizeRequest{} }
func (m *StreamingDiarizeRequest) String() string { return proto.CompactTextString(m) }
func (*StreamingDiarizeRequest) ProtoMessage()    {}
func (*StreamingDiarizeRequest) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{0}
}

func (m *StreamingDiarizeRequest) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_StreamingDiarizeRequest.Unmarshal(m, b)
}
func (m *StreamingDiarizeRequest) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_StreamingDiarizeRequest.Marshal(b, m, deterministic)
}
func (m *StreamingDiarizeRequest) XXX_Merge(src proto.Message) {
	xxx_messageInfo_StreamingDiarizeRequest.Merge(m, src)
}
func (m *StreamingDiarizeRequest) XXX_Size() int {
	return xxx_messageInfo_StreamingDiarizeRequest.Size(m)
}
func (m *StreamingDiarizeRequest) XXX_DiscardUnknown() {
	xxx_messageInfo_StreamingDiarizeRequest.DiscardUnknown(m)
}

var xxx_messageInfo_StreamingDiarizeRequest proto.InternalMessageInfo

type isStreamingDiarizeRequest_Request interface {
	isStreamingDiarizeRequest_Request()
}

type StreamingDiarizeRequest_Config struct {
	Config *DiarizationConfig `protobuf:"bytes,1,opt,name=config,proto3,oneof"`
}

type StreamingDiarizeRequest_Audio struct {
	Audio *DiarizationAudio `protobuf:"bytes,2,opt,name=audio,proto3,oneof"`
}

func (*StreamingDiarizeRequest_Config) isStreamingDiarizeRequest_Request() {}

func (*StreamingDiarizeRequest_Audio) isStreamingDiarizeRequest_Request() {}

func (m *StreamingDiarizeRequest) GetRequest() isStreamingDiarizeRequest_Request {
	if m != nil {
		return m.Request
	}
	return nil
}

func (m *StreamingDiarizeRequest) GetConfig() *DiarizationConfig {
	if x, ok := m.GetRequest().(*StreamingDiarizeRequest_Config); ok {
		return x.Config
	}
	return nil
}

func (m *StreamingDiarizeRequest) GetAudio() *DiarizationAudio {
	if x, ok := m.GetRequest().(*StreamingDiarizeRequest_Audio); ok {
		return x.Audio
	}
	return nil
}

// XXX_OneofWrappers is for the internal use of the proto package.
func (*StreamingDiarizeRequest) XXX_OneofWrappers() []interface{} {
	return []interface{}{
		(*StreamingDiarizeRequest_Config)(nil),
		(*StreamingDiarizeRequest_Audio)(nil),
	}
}

// The message sent by the server for the `Version` method.
type VersionResponse struct {
	// version of the juzu library handling the recognition.
	Juzu string `protobuf:"bytes,1,opt,name=juzu,proto3" json:"juzu,omitempty"`
	// version of the server handling these requests.
	Server               string   `protobuf:"bytes,2,opt,name=server,proto3" json:"server,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *VersionResponse) Reset()         { *m = VersionResponse{} }
func (m *VersionResponse) String() string { return proto.CompactTextString(m) }
func (*VersionResponse) ProtoMessage()    {}
func (*VersionResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{1}
}

func (m *VersionResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_VersionResponse.Unmarshal(m, b)
}
func (m *VersionResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_VersionResponse.Marshal(b, m, deterministic)
}
func (m *VersionResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_VersionResponse.Merge(m, src)
}
func (m *VersionResponse) XXX_Size() int {
	return xxx_messageInfo_VersionResponse.Size(m)
}
func (m *VersionResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_VersionResponse.DiscardUnknown(m)
}

var xxx_messageInfo_VersionResponse proto.InternalMessageInfo

func (m *VersionResponse) GetJuzu() string {
	if m != nil {
		return m.Juzu
	}
	return ""
}

func (m *VersionResponse) GetServer() string {
	if m != nil {
		return m.Server
	}
	return ""
}

// The message sent by the server for the `ListModels` method.
type ListModelsResponse struct {
	// List of models available for use that match the request.
	Models               []*Model `protobuf:"bytes,1,rep,name=models,proto3" json:"models,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ListModelsResponse) Reset()         { *m = ListModelsResponse{} }
func (m *ListModelsResponse) String() string { return proto.CompactTextString(m) }
func (*ListModelsResponse) ProtoMessage()    {}
func (*ListModelsResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{2}
}

func (m *ListModelsResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ListModelsResponse.Unmarshal(m, b)
}
func (m *ListModelsResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ListModelsResponse.Marshal(b, m, deterministic)
}
func (m *ListModelsResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ListModelsResponse.Merge(m, src)
}
func (m *ListModelsResponse) XXX_Size() int {
	return xxx_messageInfo_ListModelsResponse.Size(m)
}
func (m *ListModelsResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_ListModelsResponse.DiscardUnknown(m)
}

var xxx_messageInfo_ListModelsResponse proto.InternalMessageInfo

func (m *ListModelsResponse) GetModels() []*Model {
	if m != nil {
		return m.Models
	}
	return nil
}

// Description of a Juzu Diarization Model.
type Model struct {
	// Unique identifier of the model. This identifier is used to choose the
	// model that should be used for diarization, and is specified in the
	// `DiarizationConfig` message.
	Id string `protobuf:"bytes,1,opt,name=id,proto3" json:"id,omitempty"`
	// Model name.  This is a concise name describing the model, and maybe
	// presented to the end-user, for example, to help choose which model to use.
	Name string `protobuf:"bytes,2,opt,name=name,proto3" json:"name,omitempty"`
	// Model attributes.
	Attributes           *ModelAttributes `protobuf:"bytes,3,opt,name=attributes,proto3" json:"attributes,omitempty"`
	XXX_NoUnkeyedLiteral struct{}         `json:"-"`
	XXX_unrecognized     []byte           `json:"-"`
	XXX_sizecache        int32            `json:"-"`
}

func (m *Model) Reset()         { *m = Model{} }
func (m *Model) String() string { return proto.CompactTextString(m) }
func (*Model) ProtoMessage()    {}
func (*Model) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{3}
}

func (m *Model) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Model.Unmarshal(m, b)
}
func (m *Model) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Model.Marshal(b, m, deterministic)
}
func (m *Model) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Model.Merge(m, src)
}
func (m *Model) XXX_Size() int {
	return xxx_messageInfo_Model.Size(m)
}
func (m *Model) XXX_DiscardUnknown() {
	xxx_messageInfo_Model.DiscardUnknown(m)
}

var xxx_messageInfo_Model proto.InternalMessageInfo

func (m *Model) GetId() string {
	if m != nil {
		return m.Id
	}
	return ""
}

func (m *Model) GetName() string {
	if m != nil {
		return m.Name
	}
	return ""
}

func (m *Model) GetAttributes() *ModelAttributes {
	if m != nil {
		return m.Attributes
	}
	return nil
}

// Attributes of a Juzu Diarization Model.
type ModelAttributes struct {
	// Audio sample rate supported by the model.
	SampleRate uint32 `protobuf:"varint,1,opt,name=sample_rate,json=sampleRate,proto3" json:"sample_rate,omitempty"`
	// The type of segmentation (fixed / variable) supported by the model.
	SegmentationType     string   `protobuf:"bytes,2,opt,name=segmentation_type,json=segmentationType,proto3" json:"segmentation_type,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *ModelAttributes) Reset()         { *m = ModelAttributes{} }
func (m *ModelAttributes) String() string { return proto.CompactTextString(m) }
func (*ModelAttributes) ProtoMessage()    {}
func (*ModelAttributes) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{4}
}

func (m *ModelAttributes) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_ModelAttributes.Unmarshal(m, b)
}
func (m *ModelAttributes) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_ModelAttributes.Marshal(b, m, deterministic)
}
func (m *ModelAttributes) XXX_Merge(src proto.Message) {
	xxx_messageInfo_ModelAttributes.Merge(m, src)
}
func (m *ModelAttributes) XXX_Size() int {
	return xxx_messageInfo_ModelAttributes.Size(m)
}
func (m *ModelAttributes) XXX_DiscardUnknown() {
	xxx_messageInfo_ModelAttributes.DiscardUnknown(m)
}

var xxx_messageInfo_ModelAttributes proto.InternalMessageInfo

func (m *ModelAttributes) GetSampleRate() uint32 {
	if m != nil {
		return m.SampleRate
	}
	return 0
}

func (m *ModelAttributes) GetSegmentationType() string {
	if m != nil {
		return m.SegmentationType
	}
	return ""
}

// Collection of sequence of diarization results in a portion of audio.
// Juzu currently requires the full audio to determine which audio segments
// belong to which speaker.
type DiarizationResponse struct {
	Results              []*DiarizationResult `protobuf:"bytes,1,rep,name=results,proto3" json:"results,omitempty"`
	XXX_NoUnkeyedLiteral struct{}             `json:"-"`
	XXX_unrecognized     []byte               `json:"-"`
	XXX_sizecache        int32                `json:"-"`
}

func (m *DiarizationResponse) Reset()         { *m = DiarizationResponse{} }
func (m *DiarizationResponse) String() string { return proto.CompactTextString(m) }
func (*DiarizationResponse) ProtoMessage()    {}
func (*DiarizationResponse) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{5}
}

func (m *DiarizationResponse) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiarizationResponse.Unmarshal(m, b)
}
func (m *DiarizationResponse) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiarizationResponse.Marshal(b, m, deterministic)
}
func (m *DiarizationResponse) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiarizationResponse.Merge(m, src)
}
func (m *DiarizationResponse) XXX_Size() int {
	return xxx_messageInfo_DiarizationResponse.Size(m)
}
func (m *DiarizationResponse) XXX_DiscardUnknown() {
	xxx_messageInfo_DiarizationResponse.DiscardUnknown(m)
}

var xxx_messageInfo_DiarizationResponse proto.InternalMessageInfo

func (m *DiarizationResponse) GetResults() []*DiarizationResult {
	if m != nil {
		return m.Results
	}
	return nil
}

// Configuration for setting up a Diarizer.
type DiarizationConfig struct {
	// ID of the diarization model to use on the server.
	// Can be obtained by first getting list of models
	// on the server via ListModels().
	ModelId string `protobuf:"bytes,1,opt,name=model_id,json=modelId,proto3" json:"model_id,omitempty"`
	// The number of speakers expected in the audio;
	// If the number of speakers is unknown, set to 0.
	NumSpeakers uint32 `protobuf:"varint,2,opt,name=num_speakers,json=numSpeakers,proto3" json:"num_speakers,omitempty"`
	// Sampling rate of the audio to process.
	SampleRate uint32 `protobuf:"varint,3,opt,name=sample_rate,json=sampleRate,proto3" json:"sample_rate,omitempty"`
	// Encoding of audio data sent/streamed through the `DiarizationAudio`
	// messages.  For encodings like WAV/MP3 that have headers, the headers are
	// expected to be sent at the beginning of the stream, not in every
	// `DiarizationAudio` message.
	//
	// If not specified, the default encoding is RAW_LINEAR16.
	//
	// Depending on how they are configured, server instances of this service may
	// not support all the encodings enumerated above. They are always required to
	// accept RAW_LINEAR16.  If any other `Encoding` is specified, and it is not
	// available on the server being used, the recognition request will result in
	// an appropriate error message.
	AudioEncoding DiarizationConfig_Encoding `protobuf:"varint,4,opt,name=audio_encoding,json=audioEncoding,proto3,enum=cobaltspeech.juzu.DiarizationConfig_Encoding" json:"audio_encoding,omitempty"`
	// Unique identifier of the cubic model to be used for speech recognition. If
	// this value is specified, transcription results from the cubic model with
	// the given ID will also be returned alongside speaker labels. If it omitted
	// or blank, the results will not include transcripts, even if Cubic server
	// was included in the deployed image.
	CubicModelId string `protobuf:"bytes,5,opt,name=cubic_model_id,json=cubicModelId,proto3" json:"cubic_model_id,omitempty"`
	// Returns unformatted transcript.
	EnableRawTranscript  bool     `protobuf:"varint,6,opt,name=enable_raw_transcript,json=enableRawTranscript,proto3" json:"enable_raw_transcript,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DiarizationConfig) Reset()         { *m = DiarizationConfig{} }
func (m *DiarizationConfig) String() string { return proto.CompactTextString(m) }
func (*DiarizationConfig) ProtoMessage()    {}
func (*DiarizationConfig) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{6}
}

func (m *DiarizationConfig) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiarizationConfig.Unmarshal(m, b)
}
func (m *DiarizationConfig) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiarizationConfig.Marshal(b, m, deterministic)
}
func (m *DiarizationConfig) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiarizationConfig.Merge(m, src)
}
func (m *DiarizationConfig) XXX_Size() int {
	return xxx_messageInfo_DiarizationConfig.Size(m)
}
func (m *DiarizationConfig) XXX_DiscardUnknown() {
	xxx_messageInfo_DiarizationConfig.DiscardUnknown(m)
}

var xxx_messageInfo_DiarizationConfig proto.InternalMessageInfo

func (m *DiarizationConfig) GetModelId() string {
	if m != nil {
		return m.ModelId
	}
	return ""
}

func (m *DiarizationConfig) GetNumSpeakers() uint32 {
	if m != nil {
		return m.NumSpeakers
	}
	return 0
}

func (m *DiarizationConfig) GetSampleRate() uint32 {
	if m != nil {
		return m.SampleRate
	}
	return 0
}

func (m *DiarizationConfig) GetAudioEncoding() DiarizationConfig_Encoding {
	if m != nil {
		return m.AudioEncoding
	}
	return DiarizationConfig_RAW_LINEAR16
}

func (m *DiarizationConfig) GetCubicModelId() string {
	if m != nil {
		return m.CubicModelId
	}
	return ""
}

func (m *DiarizationConfig) GetEnableRawTranscript() bool {
	if m != nil {
		return m.EnableRawTranscript
	}
	return false
}

// Audio to be sent to the diarizer.
type DiarizationAudio struct {
	Data                 []byte   `protobuf:"bytes,1,opt,name=data,proto3" json:"data,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DiarizationAudio) Reset()         { *m = DiarizationAudio{} }
func (m *DiarizationAudio) String() string { return proto.CompactTextString(m) }
func (*DiarizationAudio) ProtoMessage()    {}
func (*DiarizationAudio) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{7}
}

func (m *DiarizationAudio) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiarizationAudio.Unmarshal(m, b)
}
func (m *DiarizationAudio) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiarizationAudio.Marshal(b, m, deterministic)
}
func (m *DiarizationAudio) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiarizationAudio.Merge(m, src)
}
func (m *DiarizationAudio) XXX_Size() int {
	return xxx_messageInfo_DiarizationAudio.Size(m)
}
func (m *DiarizationAudio) XXX_DiscardUnknown() {
	xxx_messageInfo_DiarizationAudio.DiscardUnknown(m)
}

var xxx_messageInfo_DiarizationAudio proto.InternalMessageInfo

func (m *DiarizationAudio) GetData() []byte {
	if m != nil {
		return m.Data
	}
	return nil
}

// A diarization result corresponding to a portion of audio.
type DiarizationResult struct {
	// Diarized segments containing speaker labels, timestamps and transcripts.
	Segments []*Segment `protobuf:"bytes,1,rep,name=segments,proto3" json:"segments,omitempty"`
	// Set of labels used to identify speakers in each segment.
	SpeakerLabels []string `protobuf:"bytes,2,rep,name=speaker_labels,json=speakerLabels,proto3" json:"speaker_labels,omitempty"`
	// If this is set to true, it denotes that the result is an interim partial
	// result, and could change after more audio is processed.  If unset, or set
	// to false, it denotes that this is a final result and will not change.
	//
	// Servers are not required to implement support for returning partial
	// results, and clients should generally not depend on their availability.
	IsPartial            bool     `protobuf:"varint,3,opt,name=is_partial,json=isPartial,proto3" json:"is_partial,omitempty"`
	XXX_NoUnkeyedLiteral struct{} `json:"-"`
	XXX_unrecognized     []byte   `json:"-"`
	XXX_sizecache        int32    `json:"-"`
}

func (m *DiarizationResult) Reset()         { *m = DiarizationResult{} }
func (m *DiarizationResult) String() string { return proto.CompactTextString(m) }
func (*DiarizationResult) ProtoMessage()    {}
func (*DiarizationResult) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{8}
}

func (m *DiarizationResult) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_DiarizationResult.Unmarshal(m, b)
}
func (m *DiarizationResult) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_DiarizationResult.Marshal(b, m, deterministic)
}
func (m *DiarizationResult) XXX_Merge(src proto.Message) {
	xxx_messageInfo_DiarizationResult.Merge(m, src)
}
func (m *DiarizationResult) XXX_Size() int {
	return xxx_messageInfo_DiarizationResult.Size(m)
}
func (m *DiarizationResult) XXX_DiscardUnknown() {
	xxx_messageInfo_DiarizationResult.DiscardUnknown(m)
}

var xxx_messageInfo_DiarizationResult proto.InternalMessageInfo

func (m *DiarizationResult) GetSegments() []*Segment {
	if m != nil {
		return m.Segments
	}
	return nil
}

func (m *DiarizationResult) GetSpeakerLabels() []string {
	if m != nil {
		return m.SpeakerLabels
	}
	return nil
}

func (m *DiarizationResult) GetIsPartial() bool {
	if m != nil {
		return m.IsPartial
	}
	return false
}

// A diarized segment of audio.
type Segment struct {
	// The identity of the speaker for this segment.
	SpeakerLabel string `protobuf:"bytes,1,opt,name=speaker_label,json=speakerLabel,proto3" json:"speaker_label,omitempty"`
	// Time offset relative to the beginning of audio received by the diarizer
	// and corresponding to the start of this segment.
	StartTime *duration.Duration `protobuf:"bytes,2,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// Time offset relative to the beginning of audio received by the diarizer
	// and corresponding to the end of this segment.
	EndTime *duration.Duration `protobuf:"bytes,3,opt,name=end_time,json=endTime,proto3" json:"end_time,omitempty"`
	// Text representing the transcription of the words that the speaker spoke.
	// Formatting options are set in cubicsvr.
	Transcript string `protobuf:"bytes,4,opt,name=transcript,proto3" json:"transcript,omitempty"`
	// Words in the transcript, their timestamps and confidence scores.
	Words                []*WordInfo `protobuf:"bytes,5,rep,name=words,proto3" json:"words,omitempty"`
	XXX_NoUnkeyedLiteral struct{}    `json:"-"`
	XXX_unrecognized     []byte      `json:"-"`
	XXX_sizecache        int32       `json:"-"`
}

func (m *Segment) Reset()         { *m = Segment{} }
func (m *Segment) String() string { return proto.CompactTextString(m) }
func (*Segment) ProtoMessage()    {}
func (*Segment) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{9}
}

func (m *Segment) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_Segment.Unmarshal(m, b)
}
func (m *Segment) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_Segment.Marshal(b, m, deterministic)
}
func (m *Segment) XXX_Merge(src proto.Message) {
	xxx_messageInfo_Segment.Merge(m, src)
}
func (m *Segment) XXX_Size() int {
	return xxx_messageInfo_Segment.Size(m)
}
func (m *Segment) XXX_DiscardUnknown() {
	xxx_messageInfo_Segment.DiscardUnknown(m)
}

var xxx_messageInfo_Segment proto.InternalMessageInfo

func (m *Segment) GetSpeakerLabel() string {
	if m != nil {
		return m.SpeakerLabel
	}
	return ""
}

func (m *Segment) GetStartTime() *duration.Duration {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *Segment) GetEndTime() *duration.Duration {
	if m != nil {
		return m.EndTime
	}
	return nil
}

func (m *Segment) GetTranscript() string {
	if m != nil {
		return m.Transcript
	}
	return ""
}

func (m *Segment) GetWords() []*WordInfo {
	if m != nil {
		return m.Words
	}
	return nil
}

// Word-specific information for recognized words.
type WordInfo struct {
	// The actual word in the text.
	Word string `protobuf:"bytes,1,opt,name=word,proto3" json:"word,omitempty"`
	// Confidence estimate between 0 and 1.  A higher number represents a
	// higher likelihood that the word was correctly recognized.
	Confidence float64 `protobuf:"fixed64,2,opt,name=confidence,proto3" json:"confidence,omitempty"`
	// Time offset relative to the beginning of audio received by the recognizer
	// and corresponding to the start of this spoken word.
	StartTime *duration.Duration `protobuf:"bytes,3,opt,name=start_time,json=startTime,proto3" json:"start_time,omitempty"`
	// Duration of the current word in the spoken audio.
	Duration             *duration.Duration `protobuf:"bytes,4,opt,name=duration,proto3" json:"duration,omitempty"`
	XXX_NoUnkeyedLiteral struct{}           `json:"-"`
	XXX_unrecognized     []byte             `json:"-"`
	XXX_sizecache        int32              `json:"-"`
}

func (m *WordInfo) Reset()         { *m = WordInfo{} }
func (m *WordInfo) String() string { return proto.CompactTextString(m) }
func (*WordInfo) ProtoMessage()    {}
func (*WordInfo) Descriptor() ([]byte, []int) {
	return fileDescriptor_1e313c10efe58577, []int{10}
}

func (m *WordInfo) XXX_Unmarshal(b []byte) error {
	return xxx_messageInfo_WordInfo.Unmarshal(m, b)
}
func (m *WordInfo) XXX_Marshal(b []byte, deterministic bool) ([]byte, error) {
	return xxx_messageInfo_WordInfo.Marshal(b, m, deterministic)
}
func (m *WordInfo) XXX_Merge(src proto.Message) {
	xxx_messageInfo_WordInfo.Merge(m, src)
}
func (m *WordInfo) XXX_Size() int {
	return xxx_messageInfo_WordInfo.Size(m)
}
func (m *WordInfo) XXX_DiscardUnknown() {
	xxx_messageInfo_WordInfo.DiscardUnknown(m)
}

var xxx_messageInfo_WordInfo proto.InternalMessageInfo

func (m *WordInfo) GetWord() string {
	if m != nil {
		return m.Word
	}
	return ""
}

func (m *WordInfo) GetConfidence() float64 {
	if m != nil {
		return m.Confidence
	}
	return 0
}

func (m *WordInfo) GetStartTime() *duration.Duration {
	if m != nil {
		return m.StartTime
	}
	return nil
}

func (m *WordInfo) GetDuration() *duration.Duration {
	if m != nil {
		return m.Duration
	}
	return nil
}

func init() {
	proto.RegisterEnum("cobaltspeech.juzu.DiarizationConfig_Encoding", DiarizationConfig_Encoding_name, DiarizationConfig_Encoding_value)
	proto.RegisterType((*StreamingDiarizeRequest)(nil), "cobaltspeech.juzu.StreamingDiarizeRequest")
	proto.RegisterType((*VersionResponse)(nil), "cobaltspeech.juzu.VersionResponse")
	proto.RegisterType((*ListModelsResponse)(nil), "cobaltspeech.juzu.ListModelsResponse")
	proto.RegisterType((*Model)(nil), "cobaltspeech.juzu.Model")
	proto.RegisterType((*ModelAttributes)(nil), "cobaltspeech.juzu.ModelAttributes")
	proto.RegisterType((*DiarizationResponse)(nil), "cobaltspeech.juzu.DiarizationResponse")
	proto.RegisterType((*DiarizationConfig)(nil), "cobaltspeech.juzu.DiarizationConfig")
	proto.RegisterType((*DiarizationAudio)(nil), "cobaltspeech.juzu.DiarizationAudio")
	proto.RegisterType((*DiarizationResult)(nil), "cobaltspeech.juzu.DiarizationResult")
	proto.RegisterType((*Segment)(nil), "cobaltspeech.juzu.Segment")
	proto.RegisterType((*WordInfo)(nil), "cobaltspeech.juzu.WordInfo")
}

func init() { proto.RegisterFile("juzu.proto", fileDescriptor_1e313c10efe58577) }

var fileDescriptor_1e313c10efe58577 = []byte{
	// 909 bytes of a gzipped FileDescriptorProto
	0x1f, 0x8b, 0x08, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02, 0xff, 0x94, 0x55, 0xdd, 0x6e, 0x1b, 0x45,
	0x14, 0xce, 0xfa, 0xdf, 0xc7, 0x3f, 0xb1, 0xa7, 0xd0, 0xba, 0x6e, 0x29, 0x66, 0xfb, 0x23, 0x0b,
	0x84, 0xd3, 0x1a, 0xa8, 0x90, 0x10, 0x95, 0x9c, 0x34, 0x55, 0x83, 0x52, 0x84, 0x26, 0xa1, 0x91,
	0xe0, 0x62, 0x35, 0xf6, 0x4e, 0xcc, 0x80, 0x77, 0x76, 0xd9, 0x99, 0x6d, 0x94, 0x5e, 0x21, 0xde,
	0x00, 0xf1, 0x04, 0xdc, 0x72, 0xc3, 0x8b, 0x70, 0xd7, 0x57, 0xe0, 0x05, 0x78, 0x03, 0xb4, 0x67,
	0x66, 0x17, 0xd7, 0x76, 0x13, 0xb8, 0xdb, 0x39, 0xe7, 0x7c, 0xe7, 0x9c, 0xf9, 0xce, 0xd9, 0x6f,
	0x00, 0xbe, 0x4f, 0x5e, 0x26, 0xa3, 0x28, 0x0e, 0x75, 0x48, 0xba, 0xb3, 0x70, 0xca, 0x16, 0x5a,
	0x45, 0x9c, 0xcf, 0xbe, 0x1b, 0xa5, 0x8e, 0xfe, 0xcd, 0x79, 0x18, 0xce, 0x17, 0x7c, 0x87, 0x45,
	0x62, 0x87, 0x49, 0x19, 0x6a, 0xa6, 0x45, 0x28, 0x95, 0x01, 0xf4, 0x6f, 0x59, 0x2f, 0x9e, 0xa6,
	0xc9, 0xe9, 0x8e, 0x9f, 0xc4, 0x18, 0x60, 0xfd, 0x37, 0x56, 0xfd, 0x3c, 0x88, 0xf4, 0xb9, 0x71,
	0xba, 0xbf, 0x39, 0x70, 0xed, 0x48, 0xc7, 0x9c, 0x05, 0x42, 0xce, 0x1f, 0x0b, 0x16, 0x8b, 0x97,
	0x9c, 0xf2, 0x1f, 0x13, 0xae, 0x34, 0x79, 0x04, 0x95, 0x59, 0x28, 0x4f, 0xc5, 0xbc, 0xe7, 0x0c,
	0x9c, 0x61, 0x63, 0x7c, 0x67, 0xb4, 0xd6, 0xda, 0xc8, 0x40, 0xb0, 0xdc, 0x1e, 0xc6, 0x3e, 0xdd,
	0xa2, 0x16, 0x45, 0x3e, 0x83, 0x32, 0x4b, 0x7c, 0x11, 0xf6, 0x0a, 0x08, 0xbf, 0x7d, 0x31, 0x7c,
	0x92, 0x86, 0x3e, 0xdd, 0xa2, 0x06, 0xb3, 0x5b, 0x87, 0x6a, 0x6c, 0xfa, 0x70, 0x3f, 0x87, 0xed,
	0xe7, 0x3c, 0x56, 0x22, 0x94, 0x94, 0xab, 0x28, 0x94, 0x8a, 0x13, 0x02, 0xa5, 0x14, 0x8f, 0x8d,
	0xd5, 0x29, 0x7e, 0x93, 0xab, 0x50, 0x51, 0x3c, 0x7e, 0xc1, 0x63, 0xac, 0x57, 0xa7, 0xf6, 0xe4,
	0x3e, 0x01, 0x72, 0x28, 0x94, 0x7e, 0x16, 0xfa, 0x7c, 0xa1, 0xf2, 0x0c, 0xf7, 0xa1, 0x12, 0xa0,
	0xa5, 0xe7, 0x0c, 0x8a, 0xc3, 0xc6, 0xb8, 0xb7, 0xa1, 0x3b, 0x84, 0x50, 0x1b, 0xe7, 0x86, 0x50,
	0x46, 0x03, 0x69, 0x43, 0x41, 0xf8, 0xb6, 0x74, 0x41, 0xf8, 0x69, 0x33, 0x92, 0x05, 0xdc, 0x96,
	0xc5, 0x6f, 0xb2, 0x0b, 0xc0, 0xb4, 0x8e, 0xc5, 0x34, 0xd1, 0x5c, 0xf5, 0x8a, 0x48, 0x80, 0xfb,
	0xa6, 0x12, 0x93, 0x3c, 0x92, 0x2e, 0xa1, 0x5c, 0x0f, 0xb6, 0x57, 0xdc, 0xe4, 0x5d, 0x68, 0x28,
	0x16, 0x44, 0x0b, 0xee, 0xc5, 0x4c, 0x73, 0xec, 0xa1, 0x45, 0xc1, 0x98, 0x28, 0xd3, 0x9c, 0x7c,
	0x00, 0x5d, 0xc5, 0xe7, 0x01, 0x97, 0x66, 0x47, 0x3c, 0x7d, 0x1e, 0x65, 0x8d, 0x75, 0x96, 0x1d,
	0xc7, 0xe7, 0x11, 0x77, 0xbf, 0x86, 0x2b, 0x4b, 0x03, 0xc8, 0xa9, 0x79, 0x94, 0x52, 0xaf, 0x92,
	0x85, 0xce, 0xb8, 0xb9, 0x64, 0xf0, 0x14, 0x83, 0x69, 0x06, 0x72, 0x5f, 0x15, 0xa0, 0xbb, 0xb6,
	0x17, 0xe4, 0x3a, 0xd4, 0x90, 0x48, 0x2f, 0xe7, 0xae, 0x8a, 0xe7, 0x03, 0x9f, 0xbc, 0x07, 0x4d,
	0x99, 0x04, 0x9e, 0x8a, 0x38, 0xfb, 0x81, 0xc7, 0x0a, 0xfb, 0x6d, 0xd1, 0x86, 0x4c, 0x82, 0x23,
	0x6b, 0x5a, 0xbd, 0x78, 0x71, 0xed, 0xe2, 0xc7, 0xd0, 0xc6, 0xc5, 0xf1, 0xb8, 0x9c, 0x85, 0xbe,
	0x90, 0xf3, 0x5e, 0x69, 0xe0, 0x0c, 0xdb, 0xe3, 0x0f, 0xff, 0xcb, 0xd2, 0x8e, 0xf6, 0x2d, 0x88,
	0xb6, 0x30, 0x49, 0x76, 0x24, 0x77, 0xa0, 0x3d, 0x4b, 0xa6, 0x62, 0xe6, 0xe5, 0xad, 0x97, 0xb1,
	0xf5, 0x26, 0x5a, 0x9f, 0xd9, 0xfe, 0xc7, 0xf0, 0x36, 0x97, 0x6c, 0x8a, 0xcd, 0x9d, 0x79, 0x3a,
	0x66, 0x52, 0xcd, 0x62, 0x11, 0xe9, 0x5e, 0x65, 0xe0, 0x0c, 0x6b, 0xf4, 0x8a, 0x71, 0x52, 0x76,
	0x76, 0x9c, 0xbb, 0xdc, 0x1d, 0xa8, 0xe5, 0x55, 0x3a, 0xd0, 0xa4, 0x93, 0x13, 0xef, 0xf0, 0xe0,
	0xcb, 0xfd, 0x09, 0x7d, 0xf0, 0xb0, 0xb3, 0x45, 0xaa, 0x50, 0x3c, 0x99, 0x3c, 0xef, 0x38, 0xa4,
	0x06, 0xa5, 0x27, 0x87, 0x93, 0xbd, 0x4e, 0xc1, 0xbd, 0x07, 0x9d, 0xd5, 0xbf, 0x25, 0xdd, 0x3c,
	0x9f, 0x69, 0x86, 0x7c, 0x36, 0x29, 0x7e, 0xbb, 0xbf, 0x38, 0xaf, 0xb1, 0x6f, 0x86, 0x43, 0x1e,
	0x42, 0xcd, 0x8e, 0x3f, 0x1b, 0x6a, 0x7f, 0x03, 0x31, 0x47, 0x26, 0x84, 0xe6, 0xb1, 0xe4, 0x2e,
	0xb4, 0xed, 0x58, 0xbc, 0x05, 0x9b, 0xa6, 0xbf, 0x4b, 0x61, 0x50, 0x1c, 0xd6, 0x69, 0xcb, 0x5a,
	0x0f, 0xd1, 0x48, 0xde, 0x01, 0x10, 0xca, 0x8b, 0x58, 0xac, 0x05, 0x5b, 0xe0, 0x74, 0x6a, 0xb4,
	0x2e, 0xd4, 0x57, 0xc6, 0xe0, 0xfe, 0xed, 0x40, 0xd5, 0xe6, 0x26, 0xb7, 0xa1, 0xf5, 0x5a, 0x46,
	0xbb, 0x0c, 0xcd, 0xe5, 0x84, 0xe4, 0x53, 0x00, 0xa5, 0x59, 0xac, 0x3d, 0x2d, 0xec, 0x8f, 0xd5,
	0x18, 0x5f, 0x1f, 0x19, 0x21, 0x1b, 0x65, 0x42, 0x36, 0x7a, 0x6c, 0x85, 0x8e, 0xd6, 0x31, 0xf8,
	0x58, 0x04, 0x9c, 0x7c, 0x0c, 0x35, 0x2e, 0x7d, 0x83, 0x2b, 0x5e, 0x86, 0xab, 0x72, 0xe9, 0x23,
	0xea, 0x16, 0xc0, 0xd2, 0xd8, 0x4a, 0xd8, 0xd1, 0x92, 0x85, 0x3c, 0x80, 0xf2, 0x59, 0x18, 0xfb,
	0xaa, 0x57, 0x46, 0xee, 0x6e, 0x6c, 0xe0, 0xee, 0x24, 0x8c, 0xfd, 0x03, 0x79, 0x1a, 0x52, 0x13,
	0xe9, 0xfe, 0xe1, 0x40, 0x2d, 0xb3, 0xa5, 0x83, 0x4a, 0xad, 0x99, 0x5e, 0xa5, 0xdf, 0x69, 0x4d,
	0x14, 0x4a, 0x9f, 0xcb, 0x99, 0xb9, 0xa3, 0x43, 0x97, 0x2c, 0x2b, 0x1c, 0x14, 0xff, 0x07, 0x07,
	0x9f, 0x40, 0x2d, 0x7b, 0x03, 0xf0, 0x2e, 0x17, 0xe2, 0xf2, 0xd0, 0xf1, 0x9f, 0x05, 0x28, 0x7d,
	0x91, 0x2a, 0xe9, 0xb7, 0x50, 0xb5, 0x82, 0x4b, 0xae, 0xae, 0x01, 0xf7, 0xd3, 0xd7, 0xa3, 0xbf,
	0x49, 0xcb, 0x56, 0x44, 0xda, 0x7d, 0xeb, 0xe7, 0x57, 0x7f, 0xfd, 0x5a, 0x68, 0x93, 0x26, 0x3e,
	0x5c, 0x2f, 0x6c, 0x46, 0x1f, 0xe0, 0x5f, 0x39, 0x7e, 0x63, 0xfe, 0xbb, 0x1b, 0xf2, 0xaf, 0xab,
	0xb8, 0x7b, 0x0d, 0x4b, 0x74, 0xc9, 0x36, 0x96, 0x58, 0x08, 0xa5, 0x8d, 0x58, 0x93, 0x9f, 0x1c,
	0xe8, 0xac, 0xbe, 0x6b, 0xe4, 0xfd, 0x4d, 0x2b, 0xbf, 0xf9, 0xf1, 0xeb, 0xdf, 0xbb, 0x54, 0xf3,
	0x36, 0x5d, 0xd2, 0x37, 0x49, 0x86, 0xce, 0x7d, 0x67, 0xf7, 0xe6, 0x37, 0x95, 0x14, 0x15, 0x4d,
	0x7f, 0x2f, 0x74, 0xf7, 0x30, 0xd7, 0x91, 0xc9, 0x95, 0x72, 0x3c, 0xad, 0xe0, 0x85, 0x3f, 0xfa,
	0x27, 0x00, 0x00, 0xff, 0xff, 0xf6, 0x9a, 0x44, 0x77, 0xfb, 0x07, 0x00, 0x00,
}

// Reference imports to suppress errors if they are not otherwise used.
var _ context.Context
var _ grpc.ClientConn

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion4

// JuzuClient is the client API for Juzu service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://godoc.org/google.golang.org/grpc#ClientConn.NewStream.
type JuzuClient interface {
	// Queries the Version of the Server.
	Version(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*VersionResponse, error)
	// Retrieves a list of available diarization models.
	ListModels(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ListModelsResponse, error)
	// Performs bidirectional streaming to enable on-the-go processing of
	// audio files, as well as the option to receive partial transcripts of
	// audio along with speaker IDs. This method is not truly streaming for
	// diarization yet, as results are received after specific chunks of audio
	// have been sent. This method is only available via GRPC and not via HTTP+JSON.
	// However, a web browser may use websockets to use this service.
	StreamingDiarize(ctx context.Context, opts ...grpc.CallOption) (Juzu_StreamingDiarizeClient, error)
}

type juzuClient struct {
	cc *grpc.ClientConn
}

func NewJuzuClient(cc *grpc.ClientConn) JuzuClient {
	return &juzuClient{cc}
}

func (c *juzuClient) Version(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*VersionResponse, error) {
	out := new(VersionResponse)
	err := c.cc.Invoke(ctx, "/cobaltspeech.juzu.Juzu/Version", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *juzuClient) ListModels(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*ListModelsResponse, error) {
	out := new(ListModelsResponse)
	err := c.cc.Invoke(ctx, "/cobaltspeech.juzu.Juzu/ListModels", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *juzuClient) StreamingDiarize(ctx context.Context, opts ...grpc.CallOption) (Juzu_StreamingDiarizeClient, error) {
	stream, err := c.cc.NewStream(ctx, &_Juzu_serviceDesc.Streams[0], "/cobaltspeech.juzu.Juzu/StreamingDiarize", opts...)
	if err != nil {
		return nil, err
	}
	x := &juzuStreamingDiarizeClient{stream}
	return x, nil
}

type Juzu_StreamingDiarizeClient interface {
	Send(*StreamingDiarizeRequest) error
	Recv() (*DiarizationResponse, error)
	grpc.ClientStream
}

type juzuStreamingDiarizeClient struct {
	grpc.ClientStream
}

func (x *juzuStreamingDiarizeClient) Send(m *StreamingDiarizeRequest) error {
	return x.ClientStream.SendMsg(m)
}

func (x *juzuStreamingDiarizeClient) Recv() (*DiarizationResponse, error) {
	m := new(DiarizationResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// JuzuServer is the server API for Juzu service.
type JuzuServer interface {
	// Queries the Version of the Server.
	Version(context.Context, *empty.Empty) (*VersionResponse, error)
	// Retrieves a list of available diarization models.
	ListModels(context.Context, *empty.Empty) (*ListModelsResponse, error)
	// Performs bidirectional streaming to enable on-the-go processing of
	// audio files, as well as the option to receive partial transcripts of
	// audio along with speaker IDs. This method is not truly streaming for
	// diarization yet, as results are received after specific chunks of audio
	// have been sent. This method is only available via GRPC and not via HTTP+JSON.
	// However, a web browser may use websockets to use this service.
	StreamingDiarize(Juzu_StreamingDiarizeServer) error
}

func RegisterJuzuServer(s *grpc.Server, srv JuzuServer) {
	s.RegisterService(&_Juzu_serviceDesc, srv)
}

func _Juzu_Version_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JuzuServer).Version(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cobaltspeech.juzu.Juzu/Version",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JuzuServer).Version(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Juzu_ListModels_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(JuzuServer).ListModels(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/cobaltspeech.juzu.Juzu/ListModels",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(JuzuServer).ListModels(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Juzu_StreamingDiarize_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(JuzuServer).StreamingDiarize(&juzuStreamingDiarizeServer{stream})
}

type Juzu_StreamingDiarizeServer interface {
	Send(*DiarizationResponse) error
	Recv() (*StreamingDiarizeRequest, error)
	grpc.ServerStream
}

type juzuStreamingDiarizeServer struct {
	grpc.ServerStream
}

func (x *juzuStreamingDiarizeServer) Send(m *DiarizationResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *juzuStreamingDiarizeServer) Recv() (*StreamingDiarizeRequest, error) {
	m := new(StreamingDiarizeRequest)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

var _Juzu_serviceDesc = grpc.ServiceDesc{
	ServiceName: "cobaltspeech.juzu.Juzu",
	HandlerType: (*JuzuServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Version",
			Handler:    _Juzu_Version_Handler,
		},
		{
			MethodName: "ListModels",
			Handler:    _Juzu_ListModels_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamingDiarize",
			Handler:       _Juzu_StreamingDiarize_Handler,
			ServerStreams: true,
			ClientStreams: true,
		},
	},
	Metadata: "juzu.proto",
}
